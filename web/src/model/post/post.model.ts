import type { BoardDTO } from "@/api/board.api";
import { CdnAPI } from "@/api/cdn.api";
import type { ApiResponse } from "@/api/http";
import { type PostDTO, PostAPI } from "@/api/post.api";
import type { ThreadDTO } from "@/api/thread.api";
import type { AxiosResponse, AxiosError } from "axios";

export const GetPostTimeReadable = (dateStr: string) => {
	var date: Date = new Date(dateStr);

	const getDayOfWeek = (date: Date) => {
		switch (date.getDay()) {
		case 1: return "Mon";
		case 2: return "Tue";
		case 3: return "Wed";
		case 4: return "Thu";
		case 5: return "Fri";
		case 6: return "Sat";
		case 7: return "Sun";
		default: return "???";
		}
	}

	const getDateStr = (date: Date) => {
		const d = date.getDate();
		const m = date.getMonth() + 1;
		const y = date.getFullYear();

		const dd = String(d).padStart(2, '0');
		const mm = String(m).padStart(2, '0');
		const yy = String(y).padStart(2, '0').substring(2);

		return `${dd}/${mm}/${yy}`;
	}

	const getTimeStr = (date: Date) => {
		const h = date.getHours();
		const m = date.getMinutes();
		const s = date.getSeconds();

		const hh = String(h).padStart(2, '0');
		const mm = String(m).padStart(2, '0');
		const ss = String(s).padStart(2, '0');
		return `${hh}:${mm}:${ss}`;
	}

	return `${getDateStr(date)} (${getDayOfWeek(date)})${getTimeStr(date)}`;
}


export type PostTextToken = {
	kind: "text";
	text: string;
};

export type PostLinkToken = {
	kind: "link";
	text: string;

	href: string;
	local: boolean;
	fail: boolean;
};

export type PostToken = PostTextToken | PostLinkToken;

export const ParsePostTokens = (text: string): PostToken[] => {
	return text.split(/(\s+|##\w+|\S+)/g).map(word => {
		if (word.startsWith(">>")) {
			return {
				kind: "link",
				text: word,
			} as unknown as PostTextToken;
		}
		return {
			kind: "text",
			href: "unknown",
			text: word,
			local: true,
			fail: true,
		} as unknown as PostLinkToken;
	});
}

export interface PostImageData {
	postId: number,
	expanded: boolean,
}

export interface ProcessedPost {
	tokens: PostToken[];
	links: Record<string, PostLinkToken>;
	image: PostImageData | null;
	backlinks: number[];
}

export const ProcessPost = async (post: PostDTO,
	thread: ThreadDTO,
	board: BoardDTO,
	imageCache: Record<number, PostImageData>,
	linksCache: Record<string, PostLinkToken>,
	thread_post_nums: number[],
	): Promise<ProcessedPost> => {
	let result: ProcessedPost = {
		tokens: [],
		links: {},
		image: null,
		backlinks: [],
	};

	// Create the PostImageData object, but only if the post has an image and
	// it's not already cached.
	// TODO: This "it's not already cached" may not be needed anymore because
	// duplicate images aren't allowed.
	if (post.filename && !imageCache[post.id]) {
		result.image = {
			postId: post.id,
			expanded: false,
		};
	}

	result.tokens = ParsePostTokens(post.content);
	for (let tok of result.tokens) {
		if (tok.kind == 'link') {
			// Before the proper routes are attributed to each link, add a
			// dummy '#' href for each of the links.
			tok.href = '#';
		}
	}

	for (let tok of result.tokens) {
		if (tok.kind == 'link') {
			// Cache hit.
			if (tok.text in linksCache && linksCache[tok.text]!.href != '#') {
				const refToken: PostLinkToken = linksCache[tok.text]!;
				tok.href = refToken.href;
				tok.local = refToken.local;
				tok.fail = refToken.fail;
				continue;
			}

			// Cache miss. Determine the link.

			// (1) Split text into `link_post_board` and `link_post_num`.

			let link_post_board = board.code;
			let link_post_num = 0;
			{
				const link_text = tok.text.substring(2); // First two chars are ">>"

				if (link_text[0] == '/') {
					const j = link_text.indexOf('/', 1);

					if (j > 0) {
						link_post_board = link_text.substring(1, j);
						link_post_num = Number(link_text.substring(j + 1));
					}
				} else {
					link_post_num = Number(link_text);
				}
			}

			// (2) Check if the link points to a post in this thread.

			let post_is_local = false;
			if (link_post_board == board.code) {
				for (let p of thread_post_nums) {
					if (p == link_post_num) {
						post_is_local = true;
						break;
					}
				}
			}

			// (3) Set link token.
			// Local link => use a direct href
			// Non-local link => need to check where it points to (if anywhere).

			tok.local = post_is_local;

			if (tok.local) {
				tok.href = `#p${link_post_num}`;

				if (result.backlinks.indexOf(link_post_num) == -1) {
					result.backlinks.push(link_post_num);
				}
			} else {
				await PostAPI.GetPostByNum(link_post_board, link_post_num).then((res: AxiosResponse<ApiResponse<PostDTO>>) => {
					const post: PostDTO = res.data.data!;
					tok.href = `/${link_post_board}/thread/${post.thread_num}#p${link_post_num}`;
				}).catch((err: AxiosError) => {
					tok.fail = true;
					console.error(err);
				});
			}

			// Cache the link token.
			result.links[tok.text] = tok as PostLinkToken;
		}
	}

	return result;
}
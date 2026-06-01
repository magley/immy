import type { BoardDTO } from "@/api/board.api";
import { CdnAPI } from "@/api/cdn.api";
import type { ApiResponse } from "@/api/http";
import { type PostDTO, PostAPI } from "@/api/post.api";
import type { ThreadDTO } from "@/api/thread.api";
import type { AxiosResponse, AxiosError } from "axios";
import { isWhiteSpaceLike, isWhiteSpaceSingleLine } from "typescript";

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
		case 0: return "Sun";
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

	type: 'normal' | 'greentext' | 'redtext';
};

export type PostLinkToken = {
	kind: "link";
	text: string;

	textTrimmed: string; // Untrimmed text can cause cache misses and stuff.
	href: string;
	local: boolean;
	fail: boolean;
};

export type PostToken = PostTextToken | PostLinkToken;

export const ParsePostTokens = (text: string): PostToken[] => {
	const isRepeated = (s: string, char: string): boolean => {
		return s == char.repeat(s.length);
	}

	const startsWithHowMany = (s: string, startChar: string): number => {
		let n = 0;
		for (let c of s) {
			if (c == startChar) n++;
			else break;
		}
		return n;
	}

	const isProbablyLink = (s: string): boolean => {
		// >>1234
		// >>/sometext/
		// >>/sometext/1234
		if (s.startsWith(">>>")) return false;
		if (!s.startsWith(">>")) return false;
		const ss = s.substring(2);
		if (ss == "") return false;
		if (ss.startsWith("/")) {
			let i = ss.substring(1).indexOf("/");
			if (i == -1) {
				return false;
			}
			return true;
		} else {
			return !isNaN(+ss);
		}
	}

	const isJustWhitespace = (s: string): boolean => {
		return s.trim().length == 0;
	}

	// Split text into `parts` like this:
	//
	// ['abc', ' ', '>>>text123', '      ', \n\n', '>what ', '>>/b/1234 ', '\n']

	let parts: string[] = [];
	let p: string = "";
	for (let i = 0; i < text.length; i++) {
		const ci = text[i]!;
		if (ci == '>' || ci == '\n' || ci == ' ') {
			if (p != "") {
				if (parts.length > 0 && isRepeated(parts.at(-1)!, '>')) {
					parts[parts.length - 1] += p;
				} else {
					parts.push(p);
				}
				p = "";
			}

			if (parts.length > 0 && isRepeated(parts.at(-1)!, ci)) {
				parts[parts.length - 1] += ci;
			} else {
				parts.push(ci);
			}
		} else {
			p += ci;
		}
	}
	if (p != "") {
		if (parts.length > 0 && isRepeated(parts.at(-1)!, '>')) {
			parts[parts.length - 1] += p;
		} else {
			parts.push(p);
		}
		p = "";
	}

	// Merge consecutive non-link parts. Whether a part is a non-link is not
	// determined with maximum precision.

	let parts2: {text: string, isLink: boolean}[] = [];

	for (let i = 0; i < parts.length; i++) {
		const text: string = parts[i]!;
		const isLink = isProbablyLink(text);
		let canMergeWithPrevious: boolean = false;

		if (!isLink && parts2.length > 0) {
			const previous = parts2.at(-1)!;
			if (!previous.isLink) {
				canMergeWithPrevious = true;
			}
			if (previous.isLink && isJustWhitespace(text)) {
				canMergeWithPrevious = true;
			}
		}

		if (canMergeWithPrevious) {
			parts2[parts2.length - 1]!.text += text;
		} else {
			parts2.push({ text, isLink });
		}
	}

	// Create tokens based on `parts2`.

	const res: PostToken[] = [];
	for (let i = 0; i < parts2.length; i++) {
		const {text, isLink} = parts2[i]!;

		if (isLink) {
			res.push({
				kind: "link",
				href: "unknown",
				text,
				textTrimmed: text.trim(),
				local: true,
				fail: false,
			} as unknown as PostLinkToken);
		} else {
			const quoteArrows: number = startsWithHowMany(text, '>');
			const isGreentext = quoteArrows > 0;

			res.push({
				kind: "text",
				text,
				type: isGreentext ? 'greentext' : 'normal',
			} as unknown as PostTextToken);
		}
	}

	return res;
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
			if (tok.textTrimmed in linksCache && linksCache[tok.textTrimmed]!.href != '#') {
				const refToken: PostLinkToken = linksCache[tok.textTrimmed]!;
				tok.href = refToken.href;
				tok.local = refToken.local;
				tok.fail = refToken.fail;
				continue;
			}

			// Cache miss. Determine the link.

			// (1) Split text into `link_post_board` and `link_post_num`.

			let [link_post_board, link_post_num] = SplitPostLink(tok.textTrimmed, board.code);

			// (A) Special case: no post number, just the board code.

			if (link_post_num == 0 && tok.textTrimmed == `>>/${board.code}/`) {
				tok.local = false;
				tok.href = `/${link_post_board}`;
			}
			else {
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
			}

			// Cache the link token.
			result.links[tok.textTrimmed] = tok as PostLinkToken;
		}
	}

	return result;
}

export const SplitPostLink = (postLink: string, fallbackBoardCode: string): [string, number] => {
	postLink = postLink.trim();

	if (!postLink.startsWith(">>")) {
		console.error("postLink needs to be of format >>123 or >>/b/123. Got:", postLink);
	};

	let link_post_board = fallbackBoardCode;
	let link_post_num = 0;

	const link_text = postLink.substring(2); // First two chars are ">>"

	if (link_text[0] == '/') {
		const j = link_text.indexOf('/', 1);

		if (j > 0) {
			link_post_board = link_text.substring(1, j);
			link_post_num = Number(link_text.substring(j + 1));
		}
	} else {
		link_post_num = Number(link_text);
	}

	return [link_post_board, link_post_num];
}
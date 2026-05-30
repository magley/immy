import type { BoardDTO } from "@/api/board.api";
import type { PostDTO } from "@/api/post.api";
import type { ThreadDTO } from "@/api/thread.api";

export const GetTabTitle = (title: string): string => {
	const websiteTitle = "ImmyChan";

	if (title != "") {
		title += " - ";
	}
	title += websiteTitle;

	return title;
}

export const GetTabTitleFrom = (parts: string[]): string => {
	return GetTabTitle(parts.join(" - "));
}

export const GetTabTitleForThread = (board: BoardDTO | undefined, thread: ThreadDTO | undefined, posts: PostDTO[] | undefined, unseenPostsCount: number): string => {
	let titleParts: string[] = [];

	if (board) {
		titleParts.push(`/${board.code}/`);
	}
	if (thread && posts) {
			if (thread.subject) {
				titleParts.push(thread.subject);
			} else if (posts[0]?.content) {
				titleParts.push(posts[0].content);
			}
		}
	if (board) {
		titleParts.push(`${board.name}`);
	}

	let title = GetTabTitleFrom(titleParts);

	if (unseenPostsCount > 0) {
		title = `(${unseenPostsCount})` + title;
	}

	return title;
}

export const GetTabTitleForBoard = (board: BoardDTO | undefined, isCatalog: boolean): string => {
	let titleParts: string[] = [];

	if (board) {
		titleParts.push(`/${board.code}/`);
	}

	if (board) {
		titleParts.push(`${board.name}`);
	}

	if (isCatalog) {
		titleParts.push("Catalog");
	}

	return GetTabTitleFrom(titleParts);
}
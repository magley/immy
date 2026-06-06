import type { ThreadForCatalogDTO } from "@/api/thread.api"

export interface ThreadStats {
	posts: number,
	images: number,
	posters: number,
	page: number,
}

export enum ThreadSortModeInCatalog {
	BumpOrder = 'bumpOrder',
	LastReply = 'lastReply',
	CreationDate = 'creationDate',
	ReplyCount = 'replyCount',
	ImageCount = 'imageCount',
	UserCount = 'userCount',
}

export const SortThreadsForCatalog = (threads: ThreadForCatalogDTO[], sortBy: ThreadSortModeInCatalog) => {
	const sortKeepSticky = (arr: ThreadForCatalogDTO[], sortFunc: (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => number): ThreadForCatalogDTO[] => {
		arr = arr.sort((a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => {
			if (a.thread.sticky && !b.thread.sticky) {
				return -1;
			}
			if (!a.thread.sticky && b.thread.sticky) {
				return 1;
			}
			return sortFunc(a, b);
		});
		return arr;
	}

	switch (sortBy) {
		case ThreadSortModeInCatalog.BumpOrder:
			const cmpLastBump = (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => {
				const dateA = new Date(a.stats.last_bump);
				const dateB = new Date(b.stats.last_bump);
				return -(dateA.getTime() - dateB.getTime());
			}
			threads = sortKeepSticky(threads, cmpLastBump);
			break;
		case ThreadSortModeInCatalog.LastReply:
			threads = sortKeepSticky(threads, (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.last_post.id - b.last_post.id));
			break;
		case ThreadSortModeInCatalog.CreationDate:
			threads = sortKeepSticky(threads, (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.thread.id - b.thread.id));
			break;
		case ThreadSortModeInCatalog.ReplyCount:
			threads = sortKeepSticky(threads, (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.stats.post_count - b.stats.post_count));
			break;
		case ThreadSortModeInCatalog.ImageCount:
			threads = sortKeepSticky(threads, (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.stats.image_count - b.stats.image_count));
			break;
		case ThreadSortModeInCatalog.UserCount:
			threads = sortKeepSticky(threads, (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.stats.user_count - b.stats.user_count));
			break;
	}

}
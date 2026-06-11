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

export const SortThreadsForCatalog = (threads: ThreadForCatalogDTO[], sortBy: ThreadSortModeInCatalog, board_code: string, pinnedCanonicalForm: string[]) => {
	const sortKeepStickyAndPin = (arr: ThreadForCatalogDTO[], sortFunc: (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => number): ThreadForCatalogDTO[] => {
		arr = arr.sort((a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => {
			let valA = 0;
			let valB = 0;

			if (a.thread.sticky) valA += 5;
			if (pinnedCanonicalForm.indexOf(ThreadToCanonicalForm(board_code, a.thread.post_num)) != -1) valA += 2;

			if (b.thread.sticky) valB += 5;
			if (pinnedCanonicalForm.indexOf(ThreadToCanonicalForm(board_code, b.thread.post_num)) != -1) valB += 2;

			if (valA == 0 && valB == 0) {
				return sortFunc(a, b);
			}

			if (valA > valB) return -1;
			if (valB > valA) return 1;
			return 0;
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
			threads = sortKeepStickyAndPin(threads, cmpLastBump);
			break;
		case ThreadSortModeInCatalog.LastReply:
			threads = sortKeepStickyAndPin(threads, (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.last_post.id - b.last_post.id));
			break;
		case ThreadSortModeInCatalog.CreationDate:
			threads = sortKeepStickyAndPin(threads, (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.thread.id - b.thread.id));
			break;
		case ThreadSortModeInCatalog.ReplyCount:
			threads = sortKeepStickyAndPin(threads, (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.stats.post_count - b.stats.post_count));
			break;
		case ThreadSortModeInCatalog.ImageCount:
			threads = sortKeepStickyAndPin(threads, (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.stats.image_count - b.stats.image_count));
			break;
		case ThreadSortModeInCatalog.UserCount:
			threads = sortKeepStickyAndPin(threads, (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.stats.user_count - b.stats.user_count));
			break;
	}

}

export const ThreadToCanonicalForm = (board: string, thread_num: number) => {
	return `${board}/${thread_num}`;
}

export const ThreadFromCanonicalForm = (s: string): [string, number] => {
	const parts = s.split("/");
	return [parts[0]!, +parts[1]!];
}
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
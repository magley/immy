export interface ThreadStats {
	posts: number,
	images: number,
	posters: number,
	page: number,
}

export interface ThreadNavProps {
	board_code: string,
	jump_to_id: string | null,
	jump_to_label: string | null,
	thread_stats: ThreadStats,
}
import { BoardAPI, type BoardDTO } from "@/api/board.api";
import { PostAPI, type PostDTO } from "@/api/post.api";
import { ThreadAPI, type ThreadDTO } from "@/api/thread.api";
import { ProcessPost, type PostImageData, type PostToken, type ProcessedPost } from "./post.model";
import type { AxiosResponse } from "axios";
import type { ApiResponse } from "@/api/http";

export interface PostPeekBundle {
	post: PostDTO;
	thread: ThreadDTO;
	board: BoardDTO;
	imageData: PostImageData;
	tokens: PostToken[];
}

export const GetPostPeek = async (
	boardCode: string,
	postNum: number,
	imageCache: Record<string, PostImageData>,
	peekCache: Record<string, PostPeekBundle>,
	): Promise<PostPeekBundle> => {
	const cacheKey: string = `${boardCode}/${postNum}`;

	if (cacheKey in peekCache) {
		return peekCache[cacheKey]!;
	}

	peekCache[cacheKey] = await GetPostPeekLoad(boardCode, postNum, imageCache);
	return peekCache[cacheKey]!;
}

// This function loads everything from the API, used on cache misses.
const GetPostPeekLoad = async (boardCode: string, postNum: number, imageCache: Record<string, PostImageData>): Promise<PostPeekBundle> => {
	let result: PostPeekBundle = {
		post: undefined!,
		thread: undefined!,
		board: undefined!,
		imageData: undefined!,
		tokens: [],
	};

	await BoardAPI.GetBoard(boardCode).then(async (res: AxiosResponse<ApiResponse<BoardDTO>>) => {
		result.board = res.data.data!;
		await PostAPI.GetPostByNum(boardCode, postNum).then(async (res: AxiosResponse<ApiResponse<PostDTO>>) => {
			result.post = res.data.data!;
			await ThreadAPI.GetThread(result.post.thread_id).then(async (res: AxiosResponse<ApiResponse<ThreadDTO>>) => {
				result.thread = res.data.data!;
				await ProcessPost(result.post, result.thread, result.board, imageCache, {}, []).then((res: ProcessedPost) => {
					if (res.image) {
						result.imageData = res.image;
					}
					result.tokens = res.tokens;
				});
			});
		});
	});

	return result;
}

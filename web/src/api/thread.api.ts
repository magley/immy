import { axiosInstance, type ApiResponse } from "@/api/http";
import { type CreatePostForThreadDTO, type PostDTO } from "@/api/post.api.ts";

export interface ThreadDTO {
    id: number;
    board_id: number;
    post_num: number;
    subject: string;
    locked: bool;
    sticky: bool;
}

export interface CreateThreadDTO {
    board_code: string;
    subject: string;
    locked: bool;
    sticky: bool;
    post: CreaetPostForThreadDTO;
}

export interface UpdateThreadDTO {
    locked: bool;
    sticky: bool;
}

export interface ThreadFullDTO {
    thread: ThreadDTO,
    posts: PostDTO[],
}

export class ThreadAPI {
    static async ListThreads(offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<ThreadDTO[]>>> {
        return axiosInstance.get<ApiResponse<ThreadDTO[]>>(`/threads/?offset=${offset}&limit=${limit}`);
    }

    static async ListThreadsByBoard(boardCode: string, offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<ThreadDTO[]>>> {
        return axiosInstance.get<ApiResponse<ThreadDTO[]>>(`/threads/board/${boardCode}?offset=${offset}&limit=${limit}`);
    }

    static async GetThread(threadId: number): Promise<AxiosResponse<ApiResponse<ThreadDTO>>> {
        return axiosInstance.get<ApiResponse<ThreadDTO>>(`/threads/${threadId}`);
    }

    static async GetFullThread(threadId: number): Promise<AxiosResponse<ApiResponse<ThreadFullDTO>>> {
        return axiosInstance.get<ApiResponse<ThreadDTO>>(`/threads/${threadId}/full`);
    }

    static async GetThreadByNum(boardCode: string, num: number): Promise<AxiosResponse<ApiResponse<ThreadDTO>>> {
        return axiosInstance.get<ApiResponse<ThreadDTO>>(`/threads/board/${boardCode}/${num}`);
    }

    static async GetFullThreadByNum(boardCode: string, num: number): Promise<AxiosResponse<ApiResponse<ThreadFullDTO>>> {
        return axiosInstance.get<ApiResponse<ThreadDTO>>(`/threads/board/${boardCode}/${num}/full`);
    }
    
    static async CreateThread(dto: CreateThreadDTO): Promise<AxiosResponse<ApiResponse<ThreadDTO>>> {
        return axiosInstance.post<ApiResponse<ThreadDTO>>(`/threads/`, dto);
    }

    static async UpdateThread(threadId: number, dto: UpdateThreadDTO): Promise<AxiosResponse<ApiResponse<ThreadDTO>>> {
        return axiosInstance.put<ApiResponse<ThreadDTO>>(`/threads/${threadId}`, dto);
    }

    static async DeleteThread(threadId: number): Promise<AxiosResponse<ApiResponse<number>>> {
        return axiosInstance.delete<ApiResponse<number>>(`/threads/${threadId}`);
    }
}

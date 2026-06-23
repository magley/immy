import { axiosInstance, type ApiResponse } from "@/api/http";
import { type CreatePostForThreadDTO, type PostDTO } from "@/api/post.api.ts";
import type { AxiosResponse } from "axios";

export interface ThreadDTO {
    id: number;
    deleted_at: string;
    board_id: number;
    post_num: number;
    subject: string;
    locked: boolean;
    sticky: boolean;
    archived: boolean;
    archived_at: string;
    auto_cycle: number;
}

export interface CreateThreadDTO {
    board_code: string;
    subject: string;
    locked: boolean;
    sticky: boolean;
    post: CreatePostForThreadDTO;
}

export interface UpdateThreadDTO {
    locked: boolean;
    sticky: boolean;
    auto_cycle: number;
}

export interface ThreadFullDTO {
    thread: ThreadDTO;
    posts: PostDTO[];
}

export interface ThreadForCatalogDTO {
    thread: ThreadDTO;
    post: PostDTO;
    last_post: PostDTO;
    stats: ThreadStats;
}

export interface ThreadForHomeDTO {
    thread: ThreadDTO;
    posts: PostDTO[]; // Initially it includes the OP post and N last posts.
    stats: ThreadStats;
}

export interface ThreadStats {
    post_count: number;
    image_count: number;
    user_count: number;
    last_bump: string;
}

export class ThreadAPI {
    static async ListThreads(offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<ThreadDTO[]>>> {
        return axiosInstance.get(`/threads/?offset=${offset}&limit=${limit}`);
    }

    static async ListThreadsByBoard(boardCode: string, offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<ThreadDTO[]>>> {
        return axiosInstance.get(`/threads/board/${boardCode}?offset=${offset}&limit=${limit}`);
    }

    static async GetThread(threadId: number): Promise<AxiosResponse<ApiResponse<ThreadDTO>>> {
        return axiosInstance.get(`/threads/${threadId}`);
    }

    static async GetFullThread(threadId: number): Promise<AxiosResponse<ApiResponse<ThreadFullDTO>>> {
        return axiosInstance.get(`/threads/${threadId}/full`);
    }

    static async GetThreadByNum(boardCode: string, num: number): Promise<AxiosResponse<ApiResponse<ThreadDTO>>> {
        return axiosInstance.get(`/threads/board/${boardCode}/${num}`);
    }

    static async GetFullThreadByNum(boardCode: string, num: number): Promise<AxiosResponse<ApiResponse<ThreadFullDTO>>> {
        return axiosInstance.get(`/threads/board/${boardCode}/${num}/full`);
    }

    static async GetThreadsForCatalog(boardCode: string): Promise<AxiosResponse<ApiResponse<ThreadForCatalogDTO[]>>> {
        return axiosInstance.get(`/threads/board/${boardCode}/catalog`);
    }

    static async GetThreadsForArchive(boardCode: string): Promise<AxiosResponse<ApiResponse<ThreadForCatalogDTO[]>>> {
        return axiosInstance.get(`/threads/board/${boardCode}/archive`);
    }

    static async GetThreadsForHome(boardCode: string, offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<ThreadForHomeDTO[]>>> {
        return axiosInstance.get(`/threads/board/${boardCode}/home?offset=${offset}&limit=${limit}`);
    }

    static async CreateThread(dto: CreateThreadDTO): Promise<AxiosResponse<ApiResponse<ThreadDTO>>> {
        return axiosInstance.post(`/threads/`, dto);
    }

    static async UpdateThread(threadId: number, dto: UpdateThreadDTO): Promise<AxiosResponse<ApiResponse<ThreadDTO>>> {
        return axiosInstance.put(`/threads/${threadId}`, dto);
    }

    static async DeleteThread(threadId: number): Promise<AxiosResponse<ApiResponse<number>>> {
        return axiosInstance.delete(`/threads/${threadId}`);
    }

    static async ArchiveThread(threadId: number): Promise<AxiosResponse<ApiResponse<ThreadDTO>>> {
        return axiosInstance.put(`/threads/${threadId}/archive`);
    }
}

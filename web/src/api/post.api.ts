import { axiosInstance, type ApiResponse } from "@/api/http";

export interface PostDTO {
    id: number;
    thread_id: number;
    board_id: number;
    num: number;
    name: string;
    tripcode: string;
    ipv4: string;
    created_at: string;
    sage: bool;
    content: string;
    filename: string;
    html: string;
}

export interface CreatePostForThreadDTO {
    name: string;
    content: string;
    filename: string;
    options: string;
}

export interface CreatePostDTO {
    name: string;
    content: string;
    filename: string;
    options: string;
    
    thread_id: number;
}

export interface UpdatePostDTO {
    name: string;
    tripcode: string;
    sage: bool;
    content: string;
    filename: string;
    html: string;
}

export class PostAPI {
    static async ListPosts(offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<PostDTO[]>>> {
        return axiosInstance.get<ApiResponse<PostDTO[]>>(`/posts/?offset=${offset}&limit=${limit}`);
    }

    static async GetPostByNum(boardCode: string, postNum: number): Promise<AxiosResponse<ApiResponse<PostDTO>>> {
        return axiosInstance.get<ApiResponse<PostDTO>>(`/posts/num/${boardCode}/${postNum}`);
    }

    static async GetPostsByThread(threadId: number): Promise<AxiosResponse<ApiResponse<PostDTO[]>>> {
        return axiosInstance.get<ApiResponse<PostDTO>>(`/posts/thread/${threadId}`);
    }

    static async CreatePost(dto: CreatePostDTO): Promise<AxiosResponse<ApiResponse<PostDTO>>> {
        return axiosInstance.post<ApiResponse<PostDTO>>(`/posts/`, dto);
    }

    static async UpdatePost(postId: number, dto: UpdatePostDTO): Promise<AxiosResponse<ApiResponse<PostDTO>>> {
        return axiosInstance.put<ApiResponse<PostDTO>>(`/posts/${postId}`, dto);
    }

    static async DeletePost(postId: number): Promise<AxiosResponse<ApiResponse<number>>> {
        return axiosInstance.delete<ApiResponse<number>>(`/posts/${postId}`);
    }
}

import { axiosInstance, type ApiResponse } from "@/api/http";
import type { AxiosResponse } from "axios";
import type { UserRole } from "./user.api";

export interface PostDTO {
    id: number;
    thread_id: number;
    thread_num: number;
    board_id: number;
    num: number;
    name: string;
    tripcode: string;
    ipv4: string;
    // user_id: number | undefined;
    user_role: UserRole | undefined;
    public_id: string | undefined;
    created_at: string;
    deleted_at: string;
    sage: boolean;
    capcode: boolean;
    content: string;
    filename: string;
    filesize: number;
    img_width: number;
    img_height: number;
    /** Base64 of MD5 hash. */
    md5: string;
    src_filename: string;
    spoiler: boolean;
    html: string;
}

export interface CreatePostForThreadDTO {
    name: string;
    content: string;
    filename: string;
    filebytes: string;
    options: string;
    spoiler: boolean;
}

export interface CreatePostDTO {
    name: string;
    content: string;
    filename: string | null;
    filebytes: string | null;
    options: string;
    spoiler: boolean;
    
    thread_id: number;
}

export interface UpdatePostDTO {
    name: string | null;
    tripcode: string | null;
    sage: boolean | null;
    content: string | null;
    filename: string | null;
    html: string | null;
}

export class PostAPI {
    static async ListPosts(offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<PostDTO[]>>> {
        return axiosInstance.get<ApiResponse<PostDTO[]>>(`/posts/?offset=${offset}&limit=${limit}`);
    }

    static async GetPostByNum(boardCode: string, postNum: number): Promise<AxiosResponse<ApiResponse<PostDTO>>> {
        return axiosInstance.get<ApiResponse<PostDTO>>(`/posts/num/${boardCode}/${postNum}`);
    }

    static async GetPostsByThread(threadId: number): Promise<AxiosResponse<ApiResponse<PostDTO[]>>> {
        return axiosInstance.get<ApiResponse<PostDTO[]>>(`/posts/thread/${threadId}`);
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

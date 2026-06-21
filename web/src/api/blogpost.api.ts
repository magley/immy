import { axiosInstance, type ApiResponse } from "@/api/http";
import type { AxiosResponse } from "axios";

export interface BlogpostDTO {
    id: number;
    title: string;
    html: string;
    author_id: number;
    author_name: string;
    created_at: string;
    deleted_at: string | null;
}

export interface BlogpostShortDTO {
    id: number;
    title: string;
    created_at: string;
}

export interface CreateBlogpostDTO {
    title: string;
    html: string;
}

export interface UpdateBlogpostDTO {
    html: string | null;
}

export class BlogpostAPI {
    static async ListBlogposts(offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<BlogpostDTO[]>>> {
        return axiosInstance.get(`/blogposts/?offset=${offset}&limit=${limit}`);
    }

    static async CreateBlogpost(dto: CreateBlogpostDTO): Promise<AxiosResponse<ApiResponse<BlogpostDTO>>> {
        return axiosInstance.post(`/blogposts/`, dto);
    }

    static async GetBlogpost(blogpostId: number): Promise<AxiosResponse<ApiResponse<BlogpostDTO>>> {
        return axiosInstance.get(`/blogposts/${blogpostId}`);
    }

    static async UpdateBlogpost(blogpostId: number, dto: UpdateBlogpostDTO): Promise<AxiosResponse<ApiResponse<BlogpostDTO>>> {
        return axiosInstance.put(`/blogposts/${blogpostId}`, dto);
    }

    static async DeleteBlogpost(blogpostId: number): Promise<AxiosResponse<ApiResponse<number>>> {
        return axiosInstance.delete(`/blogposts/${blogpostId}`);
    }
}

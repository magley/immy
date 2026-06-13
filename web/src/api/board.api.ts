import { axiosInstance, type ApiResponse } from "@/api/http";
import type { AxiosResponse } from "axios";

export interface BoardConfig {
    locked: boolean;
    hidden: boolean;
    max_file_size: number;
    reply_files_allowed: boolean;
    mime_types_allowed: string[];
    bump_limit: number;
    image_limit: number;
    flags_enabled: boolean;
    ids_enabled: boolean;
    code_enabled: boolean;
    math_enabled: boolean;
    max_threads: number;
    allow_spoilers: boolean;
    spoiler_image: string;
}

export interface BoardMeta {
    post_count: number;
    bytes_uploaded: number;
}

export interface BoardDTO {
    id: number;
    name: string;
    code: string;
    description: string;
    created_at: string;
    deleted_at: string;
    meta: BoardMeta;
    config: BoardConfig;
}

export interface CreateBoardDTO {
    name: string;
    code: string;
    description: string | null;
    config: BoardConfig;
}

export interface UpdateBoardDTO {
    name: string | null;
    code: string | null;
    description: string | null;
    locked: boolean | null;
    hidden: boolean | null;
}

export interface BoardStatisticsDTO {
    id: number;
    code: string;
    thread_count: number;
    post_count: number;
    bytes_uploaded: number;
}

export class BoardAPI {
    static async ListBoards(offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<BoardDTO[]>>> {
        return axiosInstance.get<ApiResponse<BoardDTO[]>>(`/boards/?offset=${offset}&limit=${limit}`);
    }

    static async CreateBoard(dto: CreateBoardDTO): Promise<AxiosResponse<ApiResponse<BoardDTO>>> {
        return axiosInstance.post<ApiResponse<BoardDTO>>(`/boards/`, dto);
    }

    static async GetBoard(boardCode: string): Promise<AxiosResponse<ApiResponse<BoardDTO>>> {
        return axiosInstance.get<ApiResponse<BoardDTO>>(`/boards/code/${boardCode}`);
    }

    static async GetBoardById(id: number): Promise<AxiosResponse<ApiResponse<BoardDTO>>> {
        return axiosInstance.get<ApiResponse<BoardDTO>>(`/boards/id/${id}`);
    }

    static async UpdateBoard(boardCode: string, dto: BoardDTO): Promise<AxiosResponse<ApiResponse<BoardDTO>>> {
        return axiosInstance.put<ApiResponse<BoardDTO>>(`/boards/${boardCode}`, dto);
    }

    static async DeleteBoard(boardCode: string): Promise<AxiosResponse<ApiResponse<number>>> {
        return axiosInstance.delete<ApiResponse<number>>(`/boards/${boardCode}`);
    }
}

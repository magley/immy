import { axiosInstance, type ApiResponse } from "@/api/http";

export interface BoardDTO {
    id: number;
    name: string;
    code: string;
    description: string;
    created_at: string;
    locked: boolean;
    hidden: boolean;
    post_count: int;
}

export interface CreateBoardDTO {
    name: string;
    code: string;
    description: string | null;
}

export interface UpdateBoardDTO {
    name: string | null;
    code: string | null;
    description: string | null;
    locked: boolean | null;
    hidden: boolean | null;
}

export class BoardAPI {
    static async ListBoards(offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<BoardDTO[]>>> {
        return axiosInstance.get<ApiResponse<BoardDTO[]>>(`/boards/?offset=${offset}&limit=${limit}`);
    }

    static async CreateBoard(dto: CreateBoardDTO): Promise<AxiosResponse<ApiResponse<BoardDTO>>> {
        return axiosInstance.post<ApiResponse<BoardDTO>>(`/boards/`, dto);
    }

    static async GetBoard(boardCode: string): Promise<AxiosResponse<ApiResponse<BoardDTO>>> {
        return axiosInstance.get<ApiResponse<BoardDTO>>(`/boards/${boardCode}`);
    }

    static async UpdateBoard(boardCode: string, dto: UpdateBoardDTO): Promise<AxiosResponse<ApiResponse<BoardDTO>>> {
        return axiosInstance.put<ApiResponse<BoardDTO>>(`/boards/${boardCode}`, dto);
    }

    static async DeleteBoard(boardCode: string): Promise<AxiosResponse<ApiResponse<number>>> {
        return axiosInstance.delete<ApiResponse<number>>(`/boards/${boardCode}`);
    }
}

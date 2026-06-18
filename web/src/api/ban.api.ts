import { axiosInstance, type ApiResponse } from "@/api/http";
import type { AxiosResponse } from "axios";

export interface BanDTO {
    id: number;
    ip_start: number;
    ip_end: number | null;
    created_at: string;
    expires_at: string | null;
    deleted_at: string | null;
    board_id: number | null;
    creator_id: number;
    reason: string;
    warning: boolean;
    seen: boolean;
}

export interface CreateBanDTO {
    ip_start: string;
    ip_end: string | null;
    expires_at: string | null;
    board_id: number | null;
    reason: string;
    warning: boolean;
}

export interface UpdateBanDTO {
    seen: boolean | null;
}

export class BanAPI {
    static async ListBans(offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<BanDTO[]>>> {
        return axiosInstance.get(`/bans/?offset=${offset}&limit=${limit}`);
    }

    static async ListBansForAdmin(offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<BanDTO[]>>> {
        return axiosInstance.get(`/bans/admin?offset=${offset}&limit=${limit}`);
    }

    static async GetMyBans(): Promise<AxiosResponse<ApiResponse<BanDTO[]>>> {
        return axiosInstance.get(`/bans/my`);
    }

    static async CreateBan(dto: CreateBanDTO): Promise<AxiosResponse<ApiResponse<BanDTO>>> {
        return axiosInstance.post(`/bans/`, dto);
    }

    static async GetBan(banId: number): Promise<AxiosResponse<ApiResponse<BanDTO>>> {
        return axiosInstance.get(`/bans/${banId}`);
    }

    static async GetBanForAdmin(banId: number): Promise<AxiosResponse<ApiResponse<BanDTO>>> {
        return axiosInstance.get(`/bans/admin/${banId}`);
    }

    static async UpdateBan(banId: number, dto: UpdateBanDTO): Promise<AxiosResponse<ApiResponse<BanDTO>>> {
        return axiosInstance.put(`/bans/${banId}`, dto);
    }

    static async DeleteBan(banId: number): Promise<AxiosResponse<ApiResponse<number>>> {
        return axiosInstance.delete(`/bans/${banId}`);
    }
}

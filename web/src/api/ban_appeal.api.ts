import { axiosInstance, type ApiResponse } from "@/api/http";
import type { AxiosResponse } from "axios";

export enum BanAppealStatus {
    Pending = "pending",
    Rejected = "rejected",
    RejectedFinal = "rejected_final",
    Approved = "approved"
}

export interface BanAppealDTO {
    id: number;
    ban_id: number;
    message: string;
    status: BanAppealStatus;
    created_at: string;
    deleted_at: string | null;
    reviewed_at: string | null;
    reviewed_by: number | null;
}

export interface CreateBanAppealDTO {
    ban_id: number;
    message: string;
}

export interface UpdateBanAppealDTO {
    status: BanAppealStatus | null;
}

export class BanAppealAPI {
    static async ListBanAppeals(offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<BanAppealDTO[]>>> {
        return axiosInstance.get(`/banappeals/?offset=${offset}&limit=${limit}`);
    }

    static async CreateBanAppeal(dto: CreateBanAppealDTO): Promise<AxiosResponse<ApiResponse<BanAppealDTO>>> {
        return axiosInstance.post(`/banappeals/`, dto);
    }

    static async GetBanAppeal(banappealId: number): Promise<AxiosResponse<ApiResponse<BanAppealDTO>>> {
        return axiosInstance.get(`/banappeals/${banappealId}`);
    }

    static async UpdateBanAppeal(banappealId: number, dto: UpdateBanAppealDTO): Promise<AxiosResponse<ApiResponse<BanAppealDTO>>> {
        return axiosInstance.put(`/banappeals/${banappealId}`, dto);
    }

    static async DeleteBanAppeal(banappealId: number): Promise<AxiosResponse<ApiResponse<number>>> {
        return axiosInstance.delete(`/banappeals/${banappealId}`);
    }

    static async GetBanAppealsOfBan(banId: number): Promise<AxiosResponse<ApiResponse<BanAppealDTO[]>>> {
        return axiosInstance.get(`/banappeals/ban/${banId}`);
    }

    static async CanAppealBan(banId: number): Promise<AxiosResponse<ApiResponse<boolean>>> {
        return axiosInstance.get(`/banappeals/ban/${banId}/can`);
    }
}

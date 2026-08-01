import { axiosInstance, type ApiResponse } from "@/api/http";
import type { AxiosResponse } from "axios";

export interface ConfigDTO {
    id: number;
    posting_enabled: boolean;
}

export class ConfigAPI {
    static async GetConfig(): Promise<AxiosResponse<ApiResponse<ConfigDTO>>> {
        return axiosInstance.get(`/config/`);
    }

    static async SetPostingEnabled(posting_enabled: boolean): Promise<AxiosResponse<ApiResponse<ConfigDTO>>> {
        return axiosInstance.put(`/config/post?enabled=${posting_enabled?1:0}`);
    }
}

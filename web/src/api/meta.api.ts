import { axiosInstance, type ApiResponse } from "@/api/http";
import type { AxiosResponse } from "axios";
import type { BoardStatisticsDTO } from "./board.api";

export class MetaAPI {
    static async GetMimeTypes(): Promise<AxiosResponse<ApiResponse<string[]>>> {
        return axiosInstance.get(`/meta/mime`);
    }

    static async GetStatistics(): Promise<AxiosResponse<ApiResponse<BoardStatisticsDTO[]>>> {
        return axiosInstance.get(`/meta/stats`);
    }

}

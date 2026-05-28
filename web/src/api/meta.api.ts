import { axiosInstance, type ApiResponse } from "@/api/http";
import type { AxiosResponse } from "axios";

export class MetaAPI {
    static async GetMimeTypes(): Promise<AxiosResponse<ApiResponse<string[]>>> {
        return axiosInstance.get(`/meta/mime`);
    }
}

import axios from "axios";

export interface ApiResponse<T> {
    success: boolean;
    data: T | null;
    error: ApiErrorInfo | null;
    meta: ApiMeta | null;
}

export interface ApiErrorInfo {
    code: string;
    message: string;
}

export interface ApiMeta {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
}

export const ENV = {
    API: "http://localhost:8080/api/v1",
    CDN: "http://localhost",
};

export const axiosInstance = axios.create({
    baseURL: ENV.API,
});

export const axiosInstanceCDN = axios.create({
    baseURL: ENV.CDN,
});


axiosInstance.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem("jwt");
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        console.error(error);
        return Promise.reject(error);
    }
);

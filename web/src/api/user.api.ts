import { axiosInstance, type ApiResponse } from "@/api/http";
import type { AxiosResponse } from "axios";

export enum UserType {
    Admin = 'admin',
    Moderator = 'moderator',
    Janitor = 'janitor',
}

export interface UserDTO {
    id: number;
    username: string;
    password: string;
    type: UserType;
    created_at: string;
}

export interface CreateUserDTO {
    username: string;
    password: string;
    type: UserType;
}

export interface UpdateUserDTO {
    username: string | null;
    type: UserType | null;
}

export interface LoginUserDTO {
    username: string;
    password: string;
}

export interface LoginResponseDTO {
    id: number;
    username: string;
    type: UserType;
    jwt: string;
}

export interface AuthorzationDTO {
    role: string | undefined;
}

export class UserAPI {
    static async ListUsers(offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<UserDTO[]>>> {
        return axiosInstance.get(`/users/?offset=${offset}&limit=${limit}`);
    }

    static async CreateUser(dto: CreateUserDTO): Promise<AxiosResponse<ApiResponse<UserDTO>>> {
        return axiosInstance.post(`/users/`, dto);
    }

    static async GetUser(userId: number): Promise<AxiosResponse<ApiResponse<UserDTO>>> {
        return axiosInstance.get(`/users/${userId}`);
    }

    static async UpdateUser(userId: number, dto: UpdateUserDTO): Promise<AxiosResponse<ApiResponse<UserDTO>>> {
        return axiosInstance.put(`/users/${userId}`, dto);
    }

    static async DeleteUser(userId: number): Promise<AxiosResponse<ApiResponse<number>>> {
        return axiosInstance.delete(`/users/${userId}`);
    }
    
    static async LoginUser(dto: LoginUserDTO): Promise<AxiosResponse<ApiResponse<LoginResponseDTO>>> {
        return axiosInstance.post(`/users/login`, dto);
    }

    static async AuthorizeUser(dto: AuthorzationDTO): Promise<AxiosResponse<ApiResponse<void>>> {
        return axiosInstance.post(`/users/authorize`, dto);
    }
}

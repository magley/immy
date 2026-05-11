import { axiosInstance, type ApiResponse } from "@/api/http";

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
    created_at: Date;
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

export class UserAPI {
    static async ListUsers(offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<UserDTO[]>>> {
        return axiosInstance.get<ApiResponse<UserDTO[]>>(`/users/?offset=${offset}&limit=${limit}`);
    }

    static async CreateUser(dto: CreateUserDTO): Promise<AxiosResponse<ApiResponse<UserDTO>>> {
        return axiosInstance.post<ApiResponse<UserDTO>>(`/users/`, dto);
    }

    static async GetUser(userId: number): Promise<AxiosResponse<ApiResponse<UserDTO>>> {
        return axiosInstance.get<ApiResponse<UserDTO>>(`/users/${userId}`);
    }

    static async UpdateUser(userId: number, dto: UpdateUserDTO): Promise<AxiosResponse<ApiResponse<UserDTO>>> {
        return axiosInstance.put<ApiResponse<UserDTO>>(`/users/${userId}`, dto);
    }

    static async DeleteUser(userId: number): Promise<AxiosResponse<ApiResponse<number>>> {
        return axiosInstance.delete<ApiResponse<number>>(`/users/${userId}`);
    }
    
    static async LoginUser(dto: LoginUserDTO): Promise<AxiosResponse<ApiResponse<LoginResponseDTO>>> {
        return axiosInstance.post<ApiResponse<number>>(`/users/login`, dto);
    }
}

import { axiosInstance, type ApiResponse } from "@/api/http";
import type { AxiosResponse } from "axios";

export interface RuleDTO {
    id: number;
    title: string;
    description: string;
    is_global: boolean;
    danger: number;
    created_at: string;
    deleted_at: string;
}

export interface RuleBoardDTO {
    rule_id: number;
    board_id: number;
    created_at: string;
    deleted_at: string;
}

export interface CreateRuleDTO {
    title: string;
    description: string;
    is_global: boolean;
    danger: number;
}

export interface CreateRuleBoardDTO {
    rule_id: number;
    board_id: number;
}

export interface UpdateRuleDTO {
    title: string | null;
    description: string | null;
    is_global: boolean | null;
    danger: number | null;
}

export class RuleAPI {
    static async ListRules(offset: number = 0, limit: number = 100): Promise<AxiosResponse<ApiResponse<RuleDTO[]>>> {
        return axiosInstance.get(`/rules/?offset=${offset}&limit=${limit}`);
    }

    static async CreateRule(dto: CreateRuleDTO): Promise<AxiosResponse<ApiResponse<RuleDTO>>> {
        return axiosInstance.post(`/rules/`, dto);
    }

    static async GetRule(ruleId: number): Promise<AxiosResponse<ApiResponse<RuleDTO>>> {
        return axiosInstance.get(`/rules/${ruleId}`);
    }

    static async UpdateRule(ruleId: number, dto: UpdateRuleDTO): Promise<AxiosResponse<ApiResponse<RuleDTO>>> {
        return axiosInstance.put(`/rules/${ruleId}`, dto);
    }

    static async DeleteRule(ruleId: number): Promise<AxiosResponse<ApiResponse<number>>> {
        return axiosInstance.delete(`/rules/${ruleId}`);
    }

    // ============================================================
    //  RULE <-> BOARD

    static async CreateRuleBoard(dto: CreateRuleDTO): Promise<AxiosResponse<ApiResponse<RuleBoardDTO>>> {
        return axiosInstance.post(`/rules/board/`, dto);
    }

    static async DeleteRuleBoard(boardId: number, ruleId: number): Promise<AxiosResponse<ApiResponse<number>>> {
        return axiosInstance.delete(`/rules/board/${boardId}/${ruleId}`);
    }

    static async ListAllRuleBoards(): Promise<AxiosResponse<ApiResponse<RuleBoardDTO[]>>> {
        return axiosInstance.get(`/rules/board/`);
    }

    static async ListAllRulesOfBoard(boardId: number): Promise<AxiosResponse<ApiResponse<RuleBoardDTO[]>>> {
        return axiosInstance.get(`/rules/board/rules/${boardId}`);
    }

    static async ListAllBoardsOfRule(ruleId: number): Promise<AxiosResponse<ApiResponse<RuleBoardDTO[]>>> {
        return axiosInstance.get(`/rules/board/boards/${ruleId}`);
    }
}

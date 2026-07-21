import type { BoardDTO } from "@/api/board.api";
import type { PostDTO } from "@/api/post.api";
import type { ThreadDTO } from "@/api/thread.api";

export enum FilterAction {
    Hide, Highlight
}

export enum FilterTarget {
    Comment,
    Filename,
    MD5,
    Username,
    Tripcode,
    ThreadSubject
}

export interface Filter {
    text: string;
    target: FilterTarget;
    boards: string[];
    enabled: boolean;
    action: FilterAction; 
}

const FILTERS_KEY = "filters";

export const LoadFilters = (): Filter[] => {
    const filtersSaved: string | null = localStorage.getItem(FILTERS_KEY);
    if (filtersSaved == null) {
        return [];
    }
    return JSON.parse(filtersSaved);
}

export const SaveFilters = (filters: Filter[]) => {
    const filtersSaved: string = JSON.stringify(filters);
    localStorage.setItem(FILTERS_KEY, filtersSaved);
}


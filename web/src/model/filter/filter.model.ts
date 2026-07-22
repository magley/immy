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

/// Returns which filter among the array (`filters`) matches the specified post.
/// If multiple matches, the first one is returned.
/// If there are no matches, the function returns `null`.
export const GetFilterMatchingPost = (board: BoardDTO, thread: ThreadDTO, post: PostDTO, filters: Filter[]): Filter | null => {
    for (let filter of filters) {
        if (IsPostFilteredBy(board, thread, post, filter)) {
            return filter;
        }
    }
    return null;
}

/// Check if the given filter matches the specified post.
export const IsPostFilteredBy = (board: BoardDTO, thread: ThreadDTO, post: PostDTO, filter: Filter): boolean => {
    let boardOk = false;
    let boardCodeSafe = board.code;
    if (boardCodeSafe.startsWith("/")) boardCodeSafe = boardCodeSafe.substring(1);
    if (boardCodeSafe.endsWith("/")) boardCodeSafe = boardCodeSafe.substring(0, boardCodeSafe.length - 1);
    
    if (filter.boards.length == 0 || (filter.boards.length == 1 && filter.boards[0] == "") || filter.boards.includes("*") || filter.boards.includes(boardCodeSafe)) {
        boardOk = true;
    }
    if (!boardOk) {
        return false;
    }

    switch (filter.target) {
        case FilterTarget.Comment: return isTextMatch(post.content, filter.text);
        case FilterTarget.Filename: return isTextMatch(post.filename, filter.text);
        case FilterTarget.MD5: console.warn("Not implemented: Filtering by MD5"); break;
        case FilterTarget.ThreadSubject: return isTextMatch(thread.subject, filter.text);
        case FilterTarget.Tripcode: return isTextMatch(post.tripcode, filter.text);
        case FilterTarget.Username: return isTextMatch(post.name, filter.text);
    }

    return false;
}

const isTextMatch = (text: string, pattern: string): boolean => {
    let isMatch = false;
    try {
        isMatch = text.match(pattern) != null;
    }
    catch {
    }
    return isMatch;
}
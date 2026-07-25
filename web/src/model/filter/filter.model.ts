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
    boards: string;
    enabled: boolean;
    action: FilterAction; 
    colorHex: string;
}

const FILTERS_KEY = "filters";
const FILTERS_ENABLED_KEY = "filters_enabled";

export const LoadFilters = (): [Filter[], boolean] => {
    const filtersSaved: string | null = localStorage.getItem(FILTERS_KEY);
    if (filtersSaved == null) {
        return [[], false];
    }
    const filters: Filter[] = JSON.parse(filtersSaved);
    const enabled = localStorage.getItem(FILTERS_ENABLED_KEY) === "true";
    return [filters, enabled];
}

export const SaveFilters = (filters: Filter[], enabled: boolean) => {
    const filtersSaved: string = JSON.stringify(filters);
    localStorage.setItem(FILTERS_KEY, filtersSaved);
    localStorage.setItem(FILTERS_ENABLED_KEY, String(enabled));
}

/// Returns which filter among the array (`filters`) matches the specified post.
/// If multiple matches, the first one is returned.
/// If there are no matches, the function returns `null`.
export const GetFilterMatchingPost = (board: BoardDTO, thread: ThreadDTO, post: PostDTO, filters: Filter[]): Filter | null => {
    for (let filter of filters) {
        if (!filter.enabled) continue;
        if (IsPostFilteredBy(board, thread, post, filter)) {
            return filter;
        }
    }
    return null;
}

/// Check if the given filter matches the specified post.
export const IsPostFilteredBy = (board: BoardDTO, thread: ThreadDTO, post: PostDTO, filter: Filter): boolean => {
    if (!filter.enabled) return false;
    let boardOk = false;

    if (filter.boards.length == 0 || (filter.boards.length == 1 && filter.boards[0] == "") || filter.boards.includes("*") || isBoardMatch(board.code, filter.boards.split(","))) {
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

const isBoardMatch = (boardCode: string, filterBoards: string[]): boolean => {
    const boardCodeSafe = getBoardCodeSafe(boardCode);
    for (let b of filterBoards) {
        if (boardCodeSafe == getBoardCodeSafe(b)) return true;
    }
    return false;
}

const getBoardCodeSafe = (code: string): string => {
    let boardCodeSafe = code;
    if (boardCodeSafe.startsWith("/")) boardCodeSafe = boardCodeSafe.substring(1);
    if (boardCodeSafe.endsWith("/")) boardCodeSafe = boardCodeSafe.substring(0, boardCodeSafe.length - 1);
    boardCodeSafe = boardCodeSafe.toLowerCase();
    return boardCodeSafe;
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
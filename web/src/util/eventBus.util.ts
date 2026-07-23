import mitt from "mitt";

export enum AppEvents {
    FiltersRefreshed = "filters-refreshed"
}

export const EventBus = mitt();
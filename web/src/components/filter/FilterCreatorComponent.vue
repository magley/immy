<script setup lang="ts">
import { onMounted, ref } from 'vue';


enum FilterAction {
    Hide, Highlight
}

enum FilterTarget {
    Comment,
    Filename,
    Username,
    Tripcode,
    ThreadSubject
}

interface Filter {
    text: string;
    targets: FilterTarget[];
    boards: string[];
    enabled: boolean;
    action: FilterAction; 
}

const filters = ref<Filter[]>([]);
const FILTERS_KEY = "filters";

const loadFilters = () => {
    const filtersSaved: string | null = localStorage.getItem(FILTERS_KEY);
    if (filtersSaved == null) {
        return;
    }
    filters.value = JSON.parse(filtersSaved);
}

const saveFilters = () => {
    const filtersSaved: string = JSON.stringify(filters.value);
    localStorage.setItem(FILTERS_KEY, filtersSaved);
}

onMounted(() => {
    loadFilters();
});

</script>

<template>
    <div>
        <div v-for="filter in filters">
            {{filter.text}} {{filter.targets}} {{filter.boards}} {{filter.enabled}} {{filter.action}} 
        </div>
        <div v-if="filters.length == 0">
            No filters defined
        </div>
    </div>
</template>

<style scoped>
</style>
<script setup lang="ts">
import { onMounted, ref } from 'vue';


enum FilterAction {
    Hide, Highlight
}

enum FilterTarget {
    Comment,
    Filename,
    MD5,
    Username,
    Tripcode,
    ThreadSubject
}

interface Filter {
    text: string;
    target: FilterTarget;
    boards: string[];
    enabled: boolean;
    action: FilterAction; 
}

const filters = ref<Filter[]>([]);
const FILTERS_KEY = "filters";

const newFilterText = ref<string>("");
const newFilterTarget = ref<FilterTarget>(FilterTarget.Comment);
const newFilterBoards = ref<string>("");
const newFilterAction = ref<FilterAction>(FilterAction.Hide);

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

const addFilter = () => {
    const newFilter: Filter = {
        text: newFilterText.value,
        target: newFilterTarget.value,
        boards: newFilterBoards.value.split(","),
        enabled: true,
        action: newFilterAction.value
    };
    filters.value.push(newFilter);
    saveFilters();
}

onMounted(() => {
    loadFilters();
    // newFilterBoards.value = "..."; // TODO: Current board
});

</script>

<template>
    <div>
        <!-- Filter list -->
        <div v-for="filter in filters">
            {{filter.text}} {{filter.target}} {{filter.boards}} {{filter.enabled}} {{filter.action}} 
        </div>
        <div v-if="filters.length == 0">
            No filters defined
        </div>

        <!-- New filter -->
        <div>
            <form>
                <label for="new-filter-text">Text:</label>
                <input id="new-filter-text" type="text" placeholder="Enter text" v-model="newFilterText" />

                <label for="new-filter-targets">Target:</label>
                <select id="new-filter-targets" v-model="newFilterTarget">
                    <option :value="FilterTarget.Comment">Post Comment</option>
                    <option :value="FilterTarget.Filename">Filename</option>
                    <option :value="FilterTarget.MD5">MD5 hash of file</option>
                    <option :value="FilterTarget.ThreadSubject">Thread Subject</option>
                    <option :value="FilterTarget.Username">Username</option>
                    <option :value="FilterTarget.Tripcode">Tripcode</option>
                </select>

                <label for="new-filter-boards">Board(s):</label>
                <input id="new-filter-boards" type="text" placeholder="Comma separated boards. Wildcard (*) supported" v-model="newFilterBoards" />
                                
                <label for="new-filter-action">Action:</label>
                <select id="new-filter-action" v-model="newFilterAction">
                    <option :value="FilterAction.Hide">Hide</option>
                    <option :value="FilterAction.Highlight">Highlight</option>
                </select>

                <button type="submit" @click.prevent="addFilter">Add</button>
            </form>
        </div>
    </div>
</template>

<style scoped>
</style>
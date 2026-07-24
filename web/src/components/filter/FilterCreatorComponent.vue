<script setup lang="ts">
import { FilterAction, type Filter, FilterTarget, SaveFilters, LoadFilters } from '@/model/filter/filter.model';
import { AppEvents, EventBus } from '@/util/eventBus.util';
import { onMounted, ref } from 'vue';


const filters = ref<Filter[]>([]);
const FILTERS_KEY = "filters";

const newFilterText = ref<string>("");
const newFilterIsRegex = ref<boolean>(false);
const newFilterTarget = ref<FilterTarget>(FilterTarget.Comment);
const newFilterColorHex = ref<string>("#FF0000");
const newFilterBoards = ref<string>("");
const newFilterAction = ref<FilterAction>(FilterAction.Hide);


const addFilter = () => {
    const newFilter: Filter = {
        text: newFilterText.value,
        target: newFilterTarget.value,
        colorHex: newFilterColorHex.value,
        boards: newFilterBoards.value.split(","),
        enabled: true,
        action: newFilterAction.value
    };
    filters.value.push(newFilter);
    saveFilters();
}

onMounted(() => {
    filters.value = LoadFilters();
    // newFilterBoards.value = "..."; // TODO: Current board
});

const removeFilter = (index: number) => {
    filters.value.splice(index, 1);
    saveFilters();
}

const saveFilters = () => {
    SaveFilters(filters.value);
    EventBus.emit(AppEvents.FiltersRefreshed)
}

</script>

<template>
    <div class="filters-root-container">
        <div class="filter-list">
            <!-- Filter list -->
            <div v-for="filter, index of filters" class="filter-definition">
                <span class="field">
                    <label :for="`filter-${index}-text`">Pattern:</label>
                    <input :id="`filter-${index}-text`" v-model="filter.text" @change="saveFilters" type="text" />
                </span>

                <span class="field">
                    <label :for="`filter-${index}-target`">Target:</label>
                    <select :id="`filter-${index}-target`" v-model="filter.target" @change="saveFilters">
                        <option :value="FilterTarget.Comment">Post Comment</option>
                        <option :value="FilterTarget.Filename">Filename</option>
                        <option :value="FilterTarget.MD5">MD5 hash of file</option>
                        <option :value="FilterTarget.ThreadSubject">Thread Subject</option>
                        <option :value="FilterTarget.Username">Username</option>
                        <option :value="FilterTarget.Tripcode">Tripcode</option>
                    </select>          
                </span>

                <span class="field">
                    <label :for="`filter-${index}-boards`">Boards:</label>
                    <input :id="`filter-${index}-boards`" v-model="filter.boards" @change="saveFilters" type="text" />          
                </span>

                <span class="field">
                    <label :for="`filter-${index}-action`">Action:</label>
                    <select :id="`filter-${index}-action`" v-model="filter.action" @change="saveFilters">
                        <option :value="FilterAction.Hide">Hide</option>
                        <option :value="FilterAction.Highlight">Highlight</option>
                    </select>                
                </span>

                <span class="field">
                    <label :for="`filter-${index}-color`">Color:</label>
                    <input :id="`filter-${index}-color`" v-model="filter.colorHex" @change="saveFilters" type="color" />
                    <input :id="`filter-${index}-color-text`" v-model="filter.colorHex" @change="saveFilters" type="text" />
                </span>

                <span class="field">
                    <label :for="`filter-${index}-enabled`">Enabled:</label>
                    <input :id="`filter-${index}-enabled`" v-model="filter.enabled" type="checkbox" @change="saveFilters" />         
                </span>            
                    
                <span class="field">
                    <button @click="removeFilter(index)" class="space-left">Delete</button>
                </span>
            </div>
            <div v-if="filters.length == 0">
                No filters defined
            </div>
        </div>

        <hr/>

        <!-- New filter -->
        <div class="filter-create">
            <form>
                <span class="field">
                    <label for="new-filter-text">Pattern:</label>
                    <input id="new-filter-text" type="text" placeholder="Enter text" v-model="newFilterText" />
                </span>

                <span class="field">
                    <label for="new-filter-targets">Target:</label>
                    <select id="new-filter-targets" v-model="newFilterTarget">
                        <option :value="FilterTarget.Comment">Post Comment</option>
                        <option :value="FilterTarget.Filename">Filename</option>
                        <option :value="FilterTarget.MD5">MD5 hash of file</option>
                        <option :value="FilterTarget.ThreadSubject">Thread Subject</option>
                        <option :value="FilterTarget.Username">Username</option>
                        <option :value="FilterTarget.Tripcode">Tripcode</option>
                    </select>
                </span>

                <span class="field">
                    <label for="new-filter-boards">Board(s):</label>
                    <input id="new-filter-boards" type="text" placeholder="Comma separated boards. Wildcard (*) supported" v-model="newFilterBoards" />
                </span>

                <span class="field">
                    <label for="new-filter-action">Action:</label>
                    <select id="new-filter-action" v-model="newFilterAction">
                        <option :value="FilterAction.Hide">Hide</option>
                        <option :value="FilterAction.Highlight">Highlight</option>
                    </select>
                </span>

                <span class="field">
                    <label for="new-filter-color">Color:</label>
                    <input id="new-filter-color" v-model="newFilterColorHex" type="color" />
                    <input id="new-filter-color-text" v-model="newFilterColorHex" type="text" />
                </span>

                <span class="field">
                    <button type="submit" @click.prevent="addFilter">Add New Filter</button>
                </span>
            </form>
        </div>
    </div>
</template>

<style scoped>
    .filters-root-container {
        width: 80%;
        margin: auto;
    }

    .filter-definition {
        margin-bottom: 0.2em;  
    }

    .field {
        margin-right: 1em;
    }

    .filter-create {
        background-color: var(--post-background-color);
        padding: 0.5em;
        text-align: center;
    }
</style>
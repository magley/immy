<script setup lang="ts">
import { FilterAction, type Filter, FilterTarget, SaveFilters, LoadFilters } from '@/model/filter/filter.model';
import { AppEvents, EventBus } from '@/util/eventBus.util';
import { onMounted, ref } from 'vue';
import { RouterLink } from 'vue-router';

const filters = ref<Filter[]>([]);
const filtersEnabled = ref<boolean>(true);

const newFilterText = ref<string>("");
const newFilterTarget = ref<FilterTarget>(FilterTarget.Comment);
const newFilterColorHex = ref<string>("#FF0000");
const newFilterBoards = ref<string>("");
const newFilterAction = ref<FilterAction>(FilterAction.Hide);

const addFilter = () => {
    const newFilter: Filter = {
        text: newFilterText.value,
        target: newFilterTarget.value,
        colorHex: newFilterColorHex.value,
        boards: newFilterBoards.value,
        enabled: true,
        action: newFilterAction.value
    };
    filters.value.push(newFilter);
    saveFilters();
}

onMounted(() => {
    [filters.value, filtersEnabled.value] = LoadFilters();
    // newFilterBoards.value = "..."; // TODO: Current board
});

const removeFilter = (index: number) => {
    filters.value.splice(index, 1);
    saveFilters();
}

const saveFilters = () => {
    SaveFilters(filters.value, filtersEnabled.value);
    EventBus.emit(AppEvents.FiltersRefreshed)
}

</script>

<template>

    <div class="filters-root-container">
        <div class="filters-container">
            <div class="">
                <span class="field">
                    <label :for="`filters-enabled`">Filters Enabled:</label>
                    <input :id="`filters-enabled`" v-model="filtersEnabled" type="checkbox" @change="saveFilters" />         
                </span>
            </div>
            <div class="filter-list">
                <!-- Filter list -->
                <div v-for="filter, index of filters" class="filter-definition">
                    <span class="field">
                        <label :for="`filter-${index}-text`">Pattern:</label>
                        <input :id="`filter-${index}-text`" v-model="filter.text" @change="saveFilters" type="text" class="width-10em" />
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
                        <label :for="`filter-${index}-boards`">
                            <abbr title='Comma separated board codes (e.g. "b, /g/, /w/ ") allowed. If no board is specified, or "*" is in the list, then the filter is global. '>
                                Board(s):</abbr>
                        </label>
                        <input :id="`filter-${index}-boards`" v-model="filter.boards" @change="saveFilters" type="text" class="width-6em" />          
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
                        <input :id="`filter-${index}-color-text`" v-model="filter.colorHex" @change="saveFilters" type="text" class="width-6em" />
                    </span>

                    <span class="field">
                        <label :for="`filter-${index}-enabled`">Enabled:</label>
                        <input :id="`filter-${index}-enabled`" v-model="filter.enabled" type="checkbox" @change="saveFilters" />         
                    </span>            
                        
                    <span class="field">
                        <button @click="removeFilter(index)" class="space-left"><img src="/icons/delete.png"/></button>
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
                        <input id="new-filter-text" type="text" placeholder="Enter text" v-model="newFilterText" class="width-10em" />
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
                        <label for="new-filter-boards">
                            <abbr title='Comma separated board codes (e.g. "b, /g/, /w/ ") allowed. If no board is specified, or "*" is in the list, then the filter is global. '>
                                Board(s):
                            </abbr>
                        </label>
                        <input id="new-filter-boards" type="text" placeholder="Comma separated boards. Wildcard (*) supported" v-model="newFilterBoards" class="width-6em" />
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
                        <input id="new-filter-color-text" v-model="newFilterColorHex" type="text" class="width-6em" />
                    </span>

                    <span class="field">
                        <button type="submit" @click.prevent="addFilter">Add New Filter</button>
                    </span>
                </form>
            </div>
        </div>
    </div>
</template>

<style scoped>
    .filters-container {
        width: 80%;
        margin: auto;
    }

    .filter-definition {
        width: 90%;
        margin: auto;
        display: block;
    }

    .field {
        margin-right: 1em;
    }

    .filter-create {
        background-color: var(--highlighted-post-background-color);
        padding: 0.5em;
        text-align: center;
    }

    .width-3em {
        width: 3em;
    }

    .width-6em {
        width: 6em;
    }

    .width-10em {
        width: 10em;
    }
</style>
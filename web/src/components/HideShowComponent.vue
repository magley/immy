<script setup lang="ts">
import { onMounted, ref } from 'vue';

interface Props {
    defaultVisible?: boolean,
    centerLabel?: boolean,
    label: string,
}

const isExpanded = ref<boolean>(true);
const props = defineProps<Props>();

onMounted(() => {
    isExpanded.value = props.defaultVisible ?? false;
})

const toggleVisible = () => {
    isExpanded.value = !isExpanded.value;
}

</script>

<template>
    <div :class="{centerLabel: centerLabel}">[<a href="#" @click="toggleVisible">{{ label }}
        <span v-if="isExpanded"> 🞁</span>
        <span v-else> 🞃</span>
    </a>]</div>
    <template v-if="isExpanded">
        <slot />
    </template>

</template>

<style>
.centerLabel {
    text-align: center;
}
</style>
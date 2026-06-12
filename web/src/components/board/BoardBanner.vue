<script setup lang="ts">
	import type { BoardDTO } from '@/api/board.api';
	import { CdnAPI } from '@/api/cdn.api';
	import { onMounted, ref } from 'vue';

	interface BoardBannerProps {
		board: BoardDTO;
	}

	const props = defineProps<BoardBannerProps>();
	const titleBanners = ref<string[]>([]);
	const titleBanner = ref<string | undefined>(undefined);

	onMounted(() => {
		CdnAPI.GetTitleBanners().then((res: string[]) => {
			titleBanners.value = res
			titleBanner.value = CdnAPI.GetTitleBanner(titleBanners.value);
		});
	})

</script>

<template>
	<div v-if="titleBanner" class="website-banner">
		<a href="/">
			<img :src="titleBanner" />
		</a>
	</div>

	<div id="title">
		<h1>/{{board.code}}/ - {{board.name}}
			<img v-if="board.config.locked" src="/icons/lock.png" title="Board locked for further posts" class="icon" />
		</h1>
		<small>{{board.description}}</small>
	</div>

</template>

<style scoped>
	.website-banner {
		display: block;
		width: 100%;
		margin: auto;
		align-content: center;
		text-align: center;
	}
</style>
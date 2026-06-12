<script setup lang="ts">
	import { CdnAPI } from '@/api/cdn.api';
	import { onMounted, ref } from 'vue';

	const boardBanners = ref<string[]>([]);
	const boardBanner = ref<string | undefined>(undefined);
	const boardCode = ref<string | undefined>(undefined);

	onMounted(() => {
		CdnAPI.GetBoardBanners().then((res: string[]) => {
			boardBanners.value = res
			boardBanner.value = CdnAPI.GetBoardBanner(boardBanners.value);
			if (boardBanner.value) {
				// Assuming it's always {URL}/{board_code}_{num}.{ext}
				boardCode.value = boardBanner.value.split("/").at(-1)!.split("_")[0]!;
			}
		});
	})

</script>

<template>
	<div v-if="boardBanner" class="board-banner">
		<a :href="`/${boardCode}`" target="_blank">
			<img :src="boardBanner" />
		</a>
	</div>

</template>

<style scoped>
	.board-banner {
		display: block;
		width: 100%;
		margin: auto;
		align-content: center;
		text-align: center;
	}
</style>
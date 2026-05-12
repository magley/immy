<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { useRoute, useRouter } from "vue-router";
	import { BoardAPI, BoardDTO, CreateBoardDTO, UpdateBoardDTO } from "@/api/board.api.ts";
	
	const board = ref<BoardDTO | null>(null);

	const route = useRoute();
	const router = useRouter();
	
	onMounted(() => {
		const board_code: string = route.params.board_code;
		loadBoard(board_code);
	});
	
	const loadBoard = (boardCode: string) => {
		BoardAPI.GetBoard(boardCode).then((res: AxiosResponse<ApiResponse<BoardDTO>>) => {
			board.value = res.data.data;
		}).catch((err: AxiosError) => {
			router.push("/");
		});
	}
</script>

<template>
	<template v-if="board">
		<h1>Catalog - /{{board.code}}/ - {{board.name}}</h1>
		<small>{{board.description}}</small>
		
		<RouterLink :to="`/${route.params.board_code}`">[Return]</RouterLink>
	</template>
</template>

<style scoped>
	
</style>
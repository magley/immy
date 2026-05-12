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
			
			const thread_id: number = route.params.thread_id;
			loadThread(thread_id);
		}).catch((err: AxiosError) => {
			router.push("/");
		});
	}
	
	const loadThread = (threadId: number) => {
		console.log(threadId);
	}
</script>

<template>
	<template v-if="board">
		Thread {{ route.params.thread_id }}
		/{{board.code}}/ - {{board.name}}
		
		<RouterLink :to="`/${route.params.board_code}`">[Return]</RouterLink>
		<RouterLink :to="`/${route.params.board_code}/catalog`">[Catalog]</RouterLink>

	</template>
</template>

<style scoped>
	
</style>
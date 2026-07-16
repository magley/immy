<script setup lang="ts">
	import { BoardAPI, type BoardDTO } from '@/api/board.api';
	import type { ApiResponse } from '@/api/http';
	import type { AxiosError, AxiosResponse } from 'axios';
	import { onMounted, ref } from 'vue';

	interface BoardListNavProps {
		isCatalog: boolean
	}

	const boards = ref<BoardDTO[]>([]);
	const props = defineProps<BoardListNavProps>();

	onMounted(() => {
		BoardAPI.GetAllBoards().then((res: AxiosResponse<ApiResponse<BoardDTO[]>>) => {
			boards.value = res.data.data!;
			boards.value = boards.value.filter((b) => !b.config.hidden);
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	});

	const getLinkForBoard = (board: BoardDTO): string => {
		if (props.isCatalog) {
			return `/${board.code}/catalog`;
		} else {
			return `/${board.code}`;
		}
	}
</script>

<template>
	[
	<span v-for="board, i in boards">
		<RouterLink :to="getLinkForBoard(board)">{{ board.code }}</RouterLink>
		<template v-if="i < boards.length - 1"> / </template>
	</span>
	]
	[ <RouterLink to="/bans">bans</RouterLink> | <RouterLink to="/blog">blog</RouterLink> ]
</template>
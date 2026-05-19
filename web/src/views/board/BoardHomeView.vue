<script setup lang="ts">
	import { BoardAPI, type BoardDTO } from "@/api/board.api.ts";
	import type { ApiResponse } from '@/api/http';
	import { ThreadAPI, type ThreadDTO } from "@/api/thread.api.ts";
	import CreateThreadForm from '@/components/thread/CreateThreadForm.vue';
	import type { AxiosError, AxiosResponse } from 'axios';
	import { onMounted, ref } from 'vue';
	import { useRoute, useRouter } from "vue-router";
	
	
	const board = ref<BoardDTO | null>(null);

	const route = useRoute();
	const router = useRouter();
	
	const threads = ref<ThreadDTO[]>([]);
	const threadsError = ref<string | undefined>(undefined);
	
	onMounted(() => {
		const board_code: string = route.params.board_code as string;
		loadBoard(board_code);
	});
	
	const loadBoard = (boardCode: string) => {
		BoardAPI.GetBoard(boardCode).then((res: AxiosResponse<ApiResponse<BoardDTO>>) => {
			board.value = res.data.data;
			loadThreads();
		}).catch((err: AxiosError) => {
			router.push("/");
		});
	}
	
	const loadThreads = () => {
		if (board.value == null) {
			return;
		}

		threadsError.value = undefined;
		ThreadAPI.ListThreadsByBoard(board.value.code).then((res : AxiosResponse<ApiResponse<ThreadDTO[]>>) => {
			threads.value = res.data.data!;
		}).catch((err: AxiosError) => {
			threadsError.value = "Could not fetch threads";
			console.error(err);
		});
	}
</script>

<template>
	<template v-if="board">
		<h1>/{{board.code}}/ - {{board.name}}</h1>
		<small>{{board.description}}</small>
		
		<RouterLink :to="`/${route.params.board_code}/catalog`">[Catalog]</RouterLink>
		
		<hr />
		
		<CreateThreadForm :board_code="board.code" :max_size_bytes="1*1024*1024" @threadCreated="loadThreads()" />
		
		<hr />
		
		<!-- Thread list -->
		<template v-if="threadsError">
			<div class="error">{{ threadsError }}</div>
		</template>
		<template v-else>
			<ul>
				<li v-for="thread in threads">
					{{ thread.id }} | {{ thread.subject }} | {{ thread.locked }} | {{ thread.sticky }}
					<RouterLink :to="`/${board.code}/thread/${thread.post_num}`">Open</RouterLink>
				</li>
			</ul>
		</template>
		
		
	</template>
</template>

<style scoped>
	.error {
		color: red;
	}
</style>
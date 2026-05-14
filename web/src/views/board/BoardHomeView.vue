<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { useRoute, useRouter } from "vue-router";
	import { BoardAPI, BoardDTO, CreateBoardDTO, UpdateBoardDTO } from "@/api/board.api.ts";
	import { ThreadAPI, ThreadDTO, CreateThreadDTO, UpdateThreadDTO } from "@/api/thread.api.ts";
	import { PostAPI, PostDTO, CreatePostForThreadDTO, CreatePostDTO, UpdatePostDTO } from "@/api/post.api.ts";
	
	
	const board = ref<BoardDTO | null>(null);

	const route = useRoute();
	const router = useRouter();
	
	const threads = ref<ThreadDTO[]>([]);
	const threadsError = ref<string | undefined>(undefined);
	
	const createThreadPostDTO = ref<CreatePostForThreadDTO>({});
	const createThreadDTO = ref<CreateThreadDTO>({});
	const createThreadError = ref<string | undefined>(undefined);
	
	onMounted(() => {
		const board_code: string = route.params.board_code;
		loadBoard(board_code);
	});
	
	const loadBoard = (boardCode: string) => {
		BoardAPI.GetBoard(boardCode).then((res: AxiosResponse<ApiResponse<BoardDTO>>) => {
			board.value = res.data.data;
			console.log(board.value);
			loadThreads();
		}).catch((err: AxiosError) => {
			router.push("/");
		});
	}
	
	const loadThreads = () => {
		threadsError.value = undefined;
		ThreadAPI.ListThreadsByBoard(board.value.code).then((res : AxiosResponse<ApiResponse<ThreadDTO[]>>) => {
			threads.value = res.data.data;
		}).catch((err: AxiosError) => {
			threadsError.value = "Could not fetch threads";
			console.error(err);
		});
	}
	
	const onSubmitCreateThread = () => {
		const realPostDTO: CreatePostForThreadDTO = {
		    name: createThreadPostDTO.value.name,
		    content: createThreadPostDTO.value.content,
		    filename: "/",
		    options: createThreadPostDTO.value.options,
		};
		const realDTO: CreateThreadDTO = {
			board_code: board.value.code,
			subject: createThreadDTO.value.subject,
			locked: false,
			sticky: false,
			post: realPostDTO,
		};
			
		ThreadAPI.CreateThread(realDTO).then(() => {
			loadThreads();
		}).catch((err: AxiosError) => {
			createThreadError.value = "Could not create thread";
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
		
		<!-- New thread -->
		<form @submit.prevent="onSubmitCreateThread">
			<input type=text placeholder="Subject" v-model="createThreadDTO.subject"/><br/>
			<input type=text placeholder="Name" v-model="createThreadPostDTO.name"/><br/>
			<input type=text placeholder="Options" v-model="createThreadPostDTO.options"/><br/>
			<textarea placeholder="Text..." v-model="createThreadPostDTO.content"/><br/>
			Files not implemented yet...
			<br/>
			<button type=submit>Create thread</button>
			
			<template v-if="createThreadError">
				<div/>
				<span class="error">{{createThreadError}}</span>
			</template>
		</form>
		
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
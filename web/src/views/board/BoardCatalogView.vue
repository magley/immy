<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { useRoute, useRouter } from "vue-router";
	import { BoardAPI, type BoardDTO } from "@/api/board.api.ts";
	import { ThreadAPI, type ThreadFullDTO, type ThreadForCatalogDTO } from '@/api/thread.api';
	import type { AxiosError, AxiosResponse } from 'axios';
	import type { ApiResponse } from '@/api/http';
	import { CdnAPI } from '@/api/cdn.api';
	import CreateThreadForm from '@/components/thread/CreateThreadForm.vue';
import BoardViewNavList from '@/components/thread/BoardViewNavList.vue';
	
	const board = ref<BoardDTO | null>(null);
	const threads = ref<ThreadForCatalogDTO[]>([]);

	const route = useRoute();
	const router = useRouter();

	const imageSizePercentage = ref<number>(100);
	const showComment = ref<boolean>(true);

	
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
		ThreadAPI.GetThreadsForCatalog(board.value!.code).then((res: AxiosResponse<ApiResponse<ThreadForCatalogDTO[]>>) => {
			threads.value = res.data.data!;

			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
			threads.value.push(threads.value[0]);
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const onSortChanged = (sortBy: string) => {
		switch (sortBy) {
		case "bumpOrder":
			// TODO: Implement bumping/saging.
			// threads.value = threads.value.sort((a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => a.thread.id - b.thread.id);
			break;
		case "lastReply":
			// TODO: Implement getting last reply date.
			// threads.value = threads.value.sort((a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => a.thread.id - b.thread.id);
			break;
		case "creationDate":
			threads.value = threads.value.sort((a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => a.thread.id - b.thread.id);
			break;
		case "replyCount":
			threads.value = threads.value.sort((a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => a.stats.post_count.id - b.stats.post_count.id);
			break;
		case "imageCount":
			threads.value = threads.value.sort((a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => a.stats.image_count.id - b.stats.image_count.id);
			break;
		case "userCount":
			threads.value = threads.value.sort((a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => a.stats.user_count.id - b.stats.user_count.id);
			break;
		}
	}

	const onImageSizeChanged = (imageSizePctg: number) => {
		imageSizePercentage.value = imageSizePctg;
	}

	const onShowCommentChanged = (show: boolean) => {
		showComment.value = show;
	}
</script>

<template>
	<template v-if="board">
		<div id="title">
			<h1>/{{board.code}}/ - {{board.name}}</h1>
			<small>{{board.description}}</small>
		</div>
		<hr />

		<CreateThreadForm id="create-thread" :board_code="board.code" :max_size_bytes="1*1024*1024" @threadCreated="loadThreads()" />
		<hr />

		<BoardViewNavList
			:board_code="board.code"
			jump_to_id="bottom"
			jump_to_label="Bottom"
			@sort-changed="onSortChanged"
			@image-size-changed="onImageSizeChanged"
			@show-comment-changed="onShowCommentChanged"
		/>

		<hr />
		<div class="catalog-grid" :style="{ gridTemplateColumns: `repeat(auto-fit, minmax(${160 * imageSizePercentage / 100}px, 1fr))`}">
			<span v-for="thread in threads" class="catalog-post">
				<RouterLink :to="`/${board.code}/thread/${thread.thread.post_num}`">
					<img
						:src="CdnAPI.GetPostImageURI(thread.post)"
						:style="{ maxWidth: (160 * imageSizePercentage / 100) + 'px', maxHeight: (140 * imageSizePercentage / 100) + 'px', }"
					>
				</RouterLink>
				<br />

				<span class="stats">
					<abbr title="Number of replies">R</abbr>: <strong>{{ thread.stats.post_count }}</strong>
					/
					<abbr title="Number of images">I</abbr>: <strong>{{ thread.stats.image_count }}</strong>
					/
					<abbr title="Number of users">U</abbr>: <strong>{{ thread.stats.user_count }}</strong>
				</span>
				<br />

				<span class="body" v-if="showComment">
					<template v-if="thread.thread.subject"><span class="subject">{{thread.thread.subject}}</span>: </template>
					<span class="content">{{ thread.post?.content }}</span>
				</span>
			</span>
		</div>

		<hr />
		<BoardViewNavList
			:board_code="board.code"
			jump_to_id="top"
			jump_to_label="Top"
			@sort-changed="onSortChanged"
			@image-size-changed="onImageSizeChanged"
			@show-comment-changed="onShowCommentChanged"
		/>

	</template>

</template>

<style scoped>
	#title {
		text-align: center;
	}

	#create-thread {
		display: block;
		text-align: center;
		width: 100%;
		margin: auto;
	}

	.catalog-grid {
		display: grid;
		/*grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));    <---  Programatically computed */
		grid-auto-rows: 320px;
		row-gap: 30px;
		column-gap: 10px;
		margin-left: 5em;
		margin-right: 5em;
		/*background-color: lightblue;*/


		.catalog-post {
			display: flex;
			flex-direction: column;
			align-items: center;
			justify-content: center;
			overflow: hidden;
			/*background-color: white;
			border: 1px solid black;*/
			padding: 5px;

			img {
			/*	max-width: 160px;      <---  Programatically computed
				max-height: 140px;     <---  Programatically computed */
				object-fit: contain;
				box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
			}

			.stats {
				cursor: help;
				text-align: center;
				align-self: center;
				font-size: small;
			}

			.body {
				text-align: center;

				.subject {
					font-weight: bold;
				}

				.content {
				}
			}
		}
	}
</style>
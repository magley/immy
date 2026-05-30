<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { useRoute, useRouter } from "vue-router";
	import { BoardAPI, type BoardDTO } from "@/api/board.api.ts";
	import { ThreadAPI, type ThreadFullDTO, type ThreadForCatalogDTO, type ThreadDTO } from '@/api/thread.api';
	import type { AxiosError, AxiosResponse } from 'axios';
	import type { ApiResponse } from '@/api/http';
	import { CdnAPI } from '@/api/cdn.api';
	import CreateThreadForm from '@/components/thread/CreateThreadForm.vue';
	import BoardViewNavList from '@/components/thread/BoardViewNavList.vue';
	import { ThreadSortModeInCatalog } from '@/model/thread/thread.model';
	import BoardListNav from '@/components/board/BoardListNav.vue';
	import { GetTabTitleForBoard } from '@/util/tab.util';

	const board = ref<BoardDTO | undefined>(undefined);
	const threads = ref<ThreadForCatalogDTO[]>([]);

	const route = useRoute();
	const router = useRouter();

	const sortBy = ref<ThreadSortModeInCatalog>(ThreadSortModeInCatalog.BumpOrder);
	const imageSize = ref<number>(100);
	const showComment = ref<boolean>(true);

	onMounted(() => {
		const board_code: string = route.params.board_code as string;
		loadBoard(board_code);
	});
	
	const loadBoard = (boardCode: string) => {
		BoardAPI.GetBoard(boardCode).then((res: AxiosResponse<ApiResponse<BoardDTO>>) => {
			board.value = res.data.data!;
			document.title = GetTabTitleForBoard(board.value, true);
			loadThreads();
		}).catch((err: AxiosError) => {
			router.push("/");
		});
	}

	const loadThreads = () => {
		ThreadAPI.GetThreadsForCatalog(board.value!.code).then((res: AxiosResponse<ApiResponse<ThreadForCatalogDTO[]>>) => {
			threads.value = res.data.data!;
			onSortChanged(sortBy.value);
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const onSortChanged = (sortCol: ThreadSortModeInCatalog) => {
		sortBy.value = sortCol;

		switch (sortBy.value) {
		case ThreadSortModeInCatalog.BumpOrder:
			const cmpLastBump = (a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => {
				const dateA = new Date(a.stats.last_bump);
				const dateB = new Date(b.stats.last_bump);
				return -(dateA.getTime() - dateB.getTime());
			}
			threads.value = threads.value.sort(cmpLastBump);
			break;
		case ThreadSortModeInCatalog.LastReply:
			threads.value = threads.value.sort((a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.last_post.id - b.last_post.id));
			break;
		case ThreadSortModeInCatalog.CreationDate:
			threads.value = threads.value.sort((a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.thread.id - b.thread.id));
			break;
		case ThreadSortModeInCatalog.ReplyCount:
			threads.value = threads.value.sort((a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.stats.post_count - b.stats.post_count));
			break;
		case ThreadSortModeInCatalog.ImageCount:
			threads.value = threads.value.sort((a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.stats.image_count - b.stats.image_count));
			break;
		case ThreadSortModeInCatalog.UserCount:
			threads.value = threads.value.sort((a: ThreadForCatalogDTO, b: ThreadForCatalogDTO) => -(a.stats.user_count - b.stats.user_count));
			break;
		}
	}

	const onImageSizeChanged = (imageSizePctg: number) => {
		imageSize.value = imageSizePctg;
	}

	const onShowCommentChanged = (show: boolean) => {
		showComment.value = show;
	}

	const getDynamicImageStyle = (thread: ThreadForCatalogDTO): any => {
		if (thread.post.img_height > thread.post.img_width) {
			return { height: (140 * imageSize.value / 100) + 'px', };
		} else {
			return { width: (160 * imageSize.value / 100) + 'px', };
		}
	}
</script>

<template>
	<BoardListNav :isCatalog=true />

	<template v-if="board">
		<div id="title">
			<h1>/{{board.code}}/ - {{board.name}}</h1>
			<small>{{board.description}}</small>
		</div>
		<hr />

		<CreateThreadForm
			id="create-thread"
			:board_code="board.code"
			:max_size_bytes="board.config.max_file_size"
			:mime_types_allowed="board.config.mime_types_allowed"
			@threadCreated="loadThreads()"
		/>
		<hr />

		<BoardViewNavList
		:board_code="board.code"
		jump_to_id="bottom"
		jump_to_label="Bottom"
		:sort-by="sortBy"
		:image-size="imageSize"
		:show-comment="showComment"
		@sort-changed="onSortChanged"
		@image-size-changed="onImageSizeChanged"
		@show-comment-changed="onShowCommentChanged"
		/>

		<hr />
		<div class="catalog-grid" :style="{ gridTemplateColumns: `repeat(auto-fit, minmax(${160 * imageSize / 100}px, 1fr))`}">
			<span v-for="thread in threads" class="catalog-post">
				<RouterLink :to="`/${board.code}/thread/${thread.thread.post_num}`">
					<img
					:src="CdnAPI.GetPostImageThumbnailURI(thread.post)"
					:style="getDynamicImageStyle(thread)"
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
		:sort-by="sortBy"
		:image-size="imageSize"
		:show-comment="showComment"
		@sort-changed="onSortChanged"
		@image-size-changed="onImageSizeChanged"
		@show-comment-changed="onShowCommentChanged"
		/>
	</template>

	<BoardListNav :isCatalog=true />
</template>

<style scoped>
	#title {
		text-align: center;
		h1 {
			color: #af0a0f;
		}
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
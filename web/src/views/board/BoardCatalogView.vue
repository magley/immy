<script setup lang="ts">
	import { ref, onMounted, useTemplateRef, nextTick } from 'vue';
	import { useRoute, useRouter } from "vue-router";
	import { BoardAPI, type BoardDTO } from "@/api/board.api.ts";
	import { ThreadAPI, type ThreadFullDTO, type ThreadForCatalogDTO, type ThreadDTO } from '@/api/thread.api';
	import type { AxiosError, AxiosResponse } from 'axios';
	import type { ApiResponse } from '@/api/http';
	import { CdnAPI } from '@/api/cdn.api';
	import BoardViewNavList from '@/components/thread/BoardViewNavList.vue';
	import { SortThreadsForCatalog, ThreadFromPinForm, ThreadSortModeInCatalog, ThreadToPinForm } from '@/model/thread/thread.model';
	import BoardListNav from '@/components/board/BoardListNav.vue';
	import { GetTabTitleForBoard } from '@/util/tab.util';
	import CreatePostForm from '@/components/post/CreatePostForm.vue';
	import { onClickOutside } from '@vueuse/core'

	const board = ref<BoardDTO | undefined>(undefined);
	const threads = ref<ThreadForCatalogDTO[]>([]);

	const route = useRoute();
	const router = useRouter();

	const sortBy = ref<ThreadSortModeInCatalog>(ThreadSortModeInCatalog.BumpOrder);
	const imageSize = ref<number>(100);
	const showComment = ref<boolean>(true);

	const pinnedThreadIDs = ref<string[]>([]);

	const modalMenuThread = ref<ThreadForCatalogDTO | undefined>(undefined);

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
			loadPinnedThreads();
			onSortChanged(sortBy.value);
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const onSortChanged = (sortCol: ThreadSortModeInCatalog) => {
		sortBy.value = sortCol;
		SortThreadsForCatalog(threads.value, sortBy.value, board.value!.code, pinnedThreadIDs.value);
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

	const pinThread = (thread_num: number) => {
		if (board.value) {
			pinnedThreadIDs.value.push(ThreadToPinForm(board.value.code, thread_num));
			pinnedThreadIDs.value = pinnedThreadIDs.value.filter((v,i,a)=>a.indexOf(v)==i); // Make unique.

			localStorage.setItem("pinned-threads", pinnedThreadIDs.value.join(","));

			onSortChanged(sortBy.value);
		}
	}

	const unpinThread = (thread_num: number) => {
		if (board.value) {
			const threadCanonical = ThreadToPinForm(board.value.code, thread_num);
			pinnedThreadIDs.value = pinnedThreadIDs.value.filter((s) => s != threadCanonical);
			localStorage.setItem("pinned-threads", pinnedThreadIDs.value.join(","));

			onSortChanged(sortBy.value);
		}
	}

	const togglePin = (thread: ThreadForCatalogDTO) => {
		if (isPinned(thread)) {
			unpinThread(thread.thread.post_num);
		} else {
			pinThread(thread.thread.post_num);
		}
	}

	const loadPinnedThreads = () => {
		let canonicalForms = (localStorage.getItem("pinned-threads") ?? "").split(",");

		// Remove any canonical forms of this thread which aren't loaded yet.
		canonicalForms = canonicalForms.filter(
			(s: string) => {
				const [board_code, thread_num] = ThreadFromPinForm(s);
				if (board_code != board.value!.code) { return true; }
				return threads.value!.findIndex((t: ThreadForCatalogDTO) => {t.thread.post_num == thread_num}) != -1;
			}
		);

		pinnedThreadIDs.value = canonicalForms;
		localStorage.setItem("pinned-threads", pinnedThreadIDs.value.join(","));
	}

	const isPinned = (thread: ThreadForCatalogDTO): boolean => {
		if (board.value == undefined) return false;
		return pinnedThreadIDs.value.indexOf(ThreadToPinForm(board.value.code, thread.post.num)) != -1;
	}

	const modalMenu = useTemplateRef("modal-menu");

	const onClickMenuArrow = async (thread: ThreadForCatalogDTO) => {
		modalMenuThread.value = thread;

		await nextTick();

		const arrow = document.getElementById(`thread-arrow-${thread.thread.post_num}`)!;
		const rect = arrow.getBoundingClientRect();

		modalMenu.value!.style.top = `${rect.top + 16}px`;
		modalMenu.value!.style.left = `${rect.left}px`;
	}

	const closeModalMenu = () => {
		modalMenuThread.value = undefined;
	}

	onClickOutside(modalMenu, event => {
		closeModalMenu();
	});
</script>

<template>
	<BoardListNav :isCatalog=true />

	<table class="modal-menu" ref="modal-menu" id="modal-menu" v-if="modalMenuThread">
		<tbody>
			<tr>
				<td>
					<a href="#" @click.prevent="">Report Thread</a>
				</td>
			</tr>
			<tr>
				<td>
					<a href="#" @click.prevent="togglePin(modalMenuThread); closeModalMenu()">
						<template v-if="isPinned(modalMenuThread)">Unpin Thread</template>
						<template v-else>Pin Thread</template>
					</a>
				</td>
			</tr>
			<tr>
				<td>
					<a href="#" @click.prevent="">Hide Thread</a>
				</td>
			</tr>
		</tbody>
	</table>

	<template v-if="board">
		<div id="title">
			<h1>/{{board.code}}/ - {{board.name}}
				<img v-if="board.config.locked" src="/icons/lock.png" title="Board locked for further posts" class="icon" />
			</h1>
			<small>{{board.description}}</small>
		</div>

		<hr />

		<div v-if="!board.config.locked">
			<CreatePostForm
			id="create-thread"
			:thread="undefined"
			:board="board"
			:max_size_bytes="board.config.max_file_size"
			:mime_types_allowed="board.config.mime_types_allowed"
			@postCreated="loadThreads()"
			/>
			<hr />
		</div>

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
					:class="{pinned: isPinned(thread)}"
					>
				</RouterLink>
				<br />

				<span class="stats">
					<img src="/icons/sticky.png" v-if="thread.thread.sticky" title="Sticky"/>
					<img src="/icons/lock.png" v-if="thread.thread.locked" title="Locked"/>
					<abbr title="Number of replies">R</abbr>: <strong>{{ thread.stats.post_count }}</strong>
					/
					<abbr title="Number of images">I</abbr>: <strong>{{ thread.stats.image_count }}</strong>
					/
					<abbr title="Number of users">U</abbr>: <strong>{{ thread.stats.user_count }}</strong>
					/
					<a href="#" class="no-underline" :id="`thread-arrow-${thread.thread.post_num}`" @click.prevent="onClickMenuArrow(thread)">▶</a>
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
	.modal-menu {
		z-index: 10;
		position: absolute;
		border: 1px solid var(--post-border-color);
		background-color: var(--post-background-color);

		td {
			border-bottom: 1px solid var(--background-color-accent);
			padding: 0.2em;
		}

		td:hover {
			background-color: var(--background-color);
		}

		a {
			text-decoration: none;
			display: block;
		}
	}

	#create-thread {
		display: block;
		text-align: center;
		width: 30%;
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

				&.pinned {
					border: 4px dashed var(--background-color-accent);
				}

				&.pinned:hover {
					border: 4px dashed var(--banner-title-color);
				}
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
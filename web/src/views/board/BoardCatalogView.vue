<script setup lang="ts">
	import { ref, onMounted, useTemplateRef, nextTick } from 'vue';
	import { useRoute, useRouter } from "vue-router";
	import { BoardAPI, type BoardDTO } from "@/api/board.api.ts";
	import { ThreadAPI, type ThreadFullDTO, type ThreadForCatalogDTO, type ThreadDTO } from '@/api/thread.api';
	import type { AxiosError, AxiosResponse } from 'axios';
	import type { ApiResponse } from '@/api/http';
	import { CdnAPI } from '@/api/cdn.api';
	import BoardViewNavList from '@/components/thread/BoardViewNavList.vue';
	import { SortThreadsForCatalog, ThreadSortModeInCatalog, ThreadToCanonicalForm, ThreadFromCanonicalForm } from '@/model/thread/thread.model';
	import BoardListNav from '@/components/board/BoardListNav.vue';
	import { GetTabTitleForBoard } from '@/util/tab.util';
	import CreatePostForm from '@/components/post/CreatePostForm.vue';
	import { onClickOutside } from '@vueuse/core'
	import BoardBanner from '@/components/board/BoardBanner.vue';
import RandomBoardImageBanner from '@/components/board/RandomBoardImageBanner.vue';

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
			loadPinnedThreads();
			loadHiddenThreads();
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

	const onClickThreadImage = (e: MouseEvent, thread: ThreadForCatalogDTO) => {
		if (e.shiftKey) {
			toggleHiddenThread(thread);
		}
		else if (e.altKey) {
			togglePin(thread);
		}
		else {
			if (board.value) {
				const destination = `/${board.value.code}/thread/${thread.thread.post_num}`;
				router.push(destination);
			}
		}
	}

	// ----------------------------------------------------------------------------------------
	// Pinned threads
	// ----------------------------------------------------------------------------------------

	const pinnedThreadIDs = ref<string[]>([]);

	const pinThread = (thread_num: number) => {
		if (board.value) {
			pinnedThreadIDs.value.push(ThreadToCanonicalForm(board.value.code, thread_num));
			pinnedThreadIDs.value = pinnedThreadIDs.value.filter((v,i,a)=>a.indexOf(v)==i); // Make unique.

			localStorage.setItem("pinned-threads", pinnedThreadIDs.value.join(","));

			onSortChanged(sortBy.value);
		}
	}

	const unpinThread = (thread_num: number) => {
		if (board.value) {
			const threadCanonical = ThreadToCanonicalForm(board.value.code, thread_num);
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
				if (s == "") return false;
				const [board_code, thread_num] = ThreadFromCanonicalForm(s);
				if (board_code != board.value!.code) { return true; }
				let shouldStay = false;
				for (let t of threads.value ?? []) {
					if (t.thread.post_num == thread_num) {
						shouldStay = true;
						break;
					}
				}
				return shouldStay;
			}
			);

		pinnedThreadIDs.value = canonicalForms;
		localStorage.setItem("pinned-threads", pinnedThreadIDs.value.join(","));
	}

	const isPinned = (thread: ThreadForCatalogDTO): boolean => {
		if (board.value == undefined) return false;
		return pinnedThreadIDs.value.indexOf(ThreadToCanonicalForm(board.value.code, thread.post.num)) != -1;
	}

	// ----------------------------------------------------------------------------------------
	// Hidden threads
	// ----------------------------------------------------------------------------------------

	// Canonical form for each thread
	const hiddenThreads = ref<string[]>([]);
	const isShowingHiddenOnly = ref<boolean>(false);

	const hideThread = (thread: ThreadForCatalogDTO) => {
		if (board.value) {
			hiddenThreads.value.push(ThreadToCanonicalForm(board.value.code, thread.thread.post_num));
			hiddenThreads.value = hiddenThreads.value.filter((v,i,a)=>a.indexOf(v)==i); // Make unique.

			localStorage.setItem("hidden-threads", hiddenThreads.value.join(","));

			onSortChanged(sortBy.value);
		}
	}

	const unhideThread = (thread: ThreadForCatalogDTO) => {
		if (board.value) {
			const threadCanonical = ThreadToCanonicalForm(board.value.code, thread.thread.post_num);
			hiddenThreads.value = hiddenThreads.value.filter((s) => s != threadCanonical);
			localStorage.setItem("hidden-threads", hiddenThreads.value.join(","));

			onSortChanged(sortBy.value);

			if (hiddenThreads.value.length == 0 && isShowingHiddenOnly.value) {
				isShowingHiddenOnly.value = false;
			}
		}
	}

	const toggleHiddenThread = (thread: ThreadForCatalogDTO) => {
		if (isHidden(thread)) {
			unhideThread(thread);
		} else {
			hideThread(thread);
		}
	}

	const isHidden = (thread: ThreadForCatalogDTO): boolean => {
		if (board.value == undefined) return false;
		return hiddenThreads.value.indexOf(ThreadToCanonicalForm(board.value.code, thread.post.num)) != -1;
	}

	const loadHiddenThreads = () => {
		let canonicalForms = (localStorage.getItem("hidden-threads") ?? "").split(",");


		// Remove any canonical forms of this thread which aren't loaded yet.
		canonicalForms = canonicalForms.filter(
			(s: string) => {
				if (s == "") return false;
				const [board_code, thread_num] = ThreadFromCanonicalForm(s);
				if (board_code != board.value!.code) { return true; }
				let shouldStay = false;
				for (let t of threads.value ?? []) {
					if (t.thread.post_num == thread_num) {
						shouldStay = true;
						break;
					}
				}
				return shouldStay;
			}
			);

		hiddenThreads.value = canonicalForms;
		localStorage.setItem("hidden-threads", hiddenThreads.value.join(","));
	}

	const onToggleHiddenThreads = () => {
		isShowingHiddenOnly.value = !isShowingHiddenOnly.value;
	}

	// ----------------------------------------------------------------------------------------
	// Modal menu
	// ----------------------------------------------------------------------------------------

	const modalMenuThread = ref<ThreadForCatalogDTO | undefined>(undefined);
	const modalMenu = useTemplateRef("modal-menu");

	const onClickMenuArrow = async (thread: ThreadForCatalogDTO) => {
		if (modalMenuThread.value == thread) {
			closeModalMenu();
		} else {
			modalMenuThread.value = thread;
		}

		await nextTick();


		if (modalMenu.value) {
			const arrow = document.getElementById(`thread-arrow-${thread.thread.post_num}`)!;
			const rect = arrow.getBoundingClientRect();

			modalMenu.value!.style.top = `${rect.top + window.scrollY + 20}px`;
			modalMenu.value!.style.left = `${rect.left + window.scrollX}px`;
		}
	}

	const closeModalMenu = () => {
		modalMenuThread.value = undefined;
	}

	onClickOutside(modalMenu, event => {
		// Delay because it conflicts with onClickMenuArrow.
		setTimeout(() => closeModalMenu(), 0.1);
	});
</script>

<template>
	<BoardListNav :isCatalog=true />

	<!-- Modal menu -->
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
						<template v-if="isPinned(modalMenuThread)"><img src="/icons/unpin.png" /> Unpin Thread</template>
						<template v-else><img src="/icons/pin.png" /> Pin Thread</template>
					</a>
				</td>
			</tr>
			<tr>
				<td>
					<a href="#" @click.prevent="toggleHiddenThread(modalMenuThread); closeModalMenu()">
						<template v-if="isHidden(modalMenuThread)"><img src="/icons/visible.png" /> Unhide Thread</template>
						<template v-else><img src="/icons/invisible.png" /> Hide Thread</template>
					</a>
				</td>
			</tr>
		</tbody>
	</table>

	<template v-if="board">
		<BoardBanner :board="board" />

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

			<RandomBoardImageBanner />

			<hr />
		</div>

		<BoardViewNavList
		:board_code="board.code"
		jump_to_id="bottom"
		jump_to_label="Bottom"
		:sort-by="sortBy"
		:image-size="imageSize"
		:show-comment="showComment"
		:hidden-threads="hiddenThreads"
		:showing-hidden="isShowingHiddenOnly"
		@sort-changed="onSortChanged"
		@image-size-changed="onImageSizeChanged"
		@show-comment-changed="onShowCommentChanged"
		@onToggleHiddenThreads="onToggleHiddenThreads"
		/>

		<hr />
		<div class="catalog-grid" :style="{ gridTemplateColumns: `repeat(auto-fit, minmax(${160 * imageSize / 100}px, 1fr))`}">
			<template v-for="thread in threads" >
				<span v-if='isShowingHiddenOnly == isHidden(thread)' class="catalog-post">
					<div class="image-container">
						<a href="#" @click.prevent="(e) => onClickThreadImage(e, thread)">
<!-- 							<img
							:src="CdnAPI.GetPostImageThumbnailURI(thread.post)"
							:style="getDynamicImageStyle(thread)"
							:class="{pinned: isPinned(thread)}"
							>
 -->
							<!-- Spoiler -->
							<img v-if="board.config.allow_spoilers && thread.post.spoiler"
								:class="{pinned: isPinned(thread)}"
								class="spoiler"
								:style="getDynamicImageStyle(thread)"
								:src="CdnAPI.GetSpoilerURI(board.config.spoiler_image)"
							>
							<!-- Regular thumbnail -->
							<img v-else
								:class="{pinned: isPinned(thread)}"
								:style="getDynamicImageStyle(thread)"
								:src="CdnAPI.GetPostImageThumbnailURI(thread.post)"
							>
						</a>
						<div class="inside-image">
							<img src="/icons/sticky.png" v-if="thread.thread.sticky" title="Sticky"/>
							<img src="/icons/lock.png" v-if="thread.thread.locked" title="Locked"/>
							<a v-if="isPinned(thread)" href="#" @click.prevent="unpinThread(thread.thread.post_num)">
								<img src="/icons/pin.png" title="Pinned - click to unpin" />
							</a>
							<a v-if="isHidden(thread)" href="#" @click.prevent="unhideThread(thread)">
								<img src="/icons/visible.png" title="Hidden - click to unhide" />
							</a>
						</div>
					</div>

					<br />

					<span class="stats">
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
			</template>
		</div>

		<hr />
		<BoardViewNavList
		:board_code="board.code"
		jump_to_id="top"
		jump_to_label="Top"
		:sort-by="sortBy"
		:image-size="imageSize"
		:show-comment="showComment"
		:hidden-threads="hiddenThreads"
		:showing-hidden="isShowingHiddenOnly"
		@sort-changed="onSortChanged"
		@image-size-changed="onImageSizeChanged"
		@show-comment-changed="onShowCommentChanged"
		@onToggleHiddenThreads="onToggleHiddenThreads"
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


			.image-container {
				position: relative;

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

					&.spoiler {
						width: 160px !important;
						height: 160px !important;
					}
				}

				.inside-image {
					position: absolute;
					top: 0px;
					left: 0px;
					z-index: 1000;
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
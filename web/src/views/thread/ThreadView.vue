<script setup lang="ts">
	import { ref, onMounted, watch, useTemplateRef, onUnmounted } from 'vue';
	import { useRoute, useRouter } from "vue-router";
	import { BoardAPI, type BoardDTO } from "@/api/board.api.ts";
	import { ThreadAPI, type ThreadDTO, type ThreadFullDTO, type ThreadForCatalogDTO, type UpdateThreadDTO } from "@/api/thread.api.ts";
	import { type PostDTO } from "@/api/post.api.ts";
	import ThreadViewNavList from "@/components/thread/ThreadViewNavList.vue";
	import { ThreadSortModeInCatalog, type ThreadStats, SortThreadsForCatalog } from "@/model/thread/thread.model.ts";
	import CreatePostForm from '@/components/post/CreatePostForm.vue';
	import type { AxiosError, AxiosResponse } from 'axios';
	import type { ApiResponse } from '@/api/http';
	import { type PostImageData, type PostToken, type PostLinkToken, type ProcessedPost, ProcessPost, SplitPostLink } from '@/model/post/post.model';
	import { vElementVisibility } from '@vueuse/components';
	import PostComponent from '@/components/post/PostComponent.vue';
	import { AddRangeNoDuplicates } from '@/util/various.util';
	import BoardListNav from '@/components/board/BoardListNav.vue';
	import { GetPostPeek, type PostPeekBundle } from '@/model/post/post.peek';
	import { GetTabTitleForThread } from '@/util/tab.util';
	import { useTextSelection, useDraggable } from '@vueuse/core';
	import { UserRole } from '@/api/user.api';
	import BoardBanner from '@/components/board/BoardBanner.vue';
	import RandomBoardImageBanner from '@/components/board/RandomBoardImageBanner.vue';
	import GalleryMode from '@/components/thread/GalleryMode.vue';

	const route = useRoute();
	const router = useRouter();

	const board = ref<BoardDTO | undefined>(undefined);
	const thread = ref<ThreadDTO | undefined>(undefined);
	const thread_stats = ref<ThreadStats>({
		posts: 0,
		images: 0,
		posters: 0,
		page: 0
	});
	const posts = ref<PostDTO[]>([]);
	/** `post.num` of the post that should be highlighted */
	const highlightedPost = ref<number | undefined>(undefined);
	/** `post.user_id` which should be highlighted */
	const highlightedPublicIDs = ref<Record<string, boolean>>({});
	/** `post.num` => list of post numbers that link to this post */
	const backLinks = ref<Record<number, number[]>>({});
	/** `post.id` => information about the image attached to the post */
	const imageData = ref<Record<number, PostImageData>>({});
	/** `post.id` => list of tokens that make up the post's content */
	const postTokens = ref<Record<number, PostToken[]>>({});
	/** `post link text from a Token` => `PostLinkToken` from which you can
	  * get all the needed info. All other tokens with the same post link
	  * text will clone the `PostLinkToken` from here. The exact same token
	  * can be found in the `postTokens` array for that particular post. */
	const postLinks = ref<Record<string, PostLinkToken>>({});
	/** When using auto-update, the last seen post (before auto update)
	  * will set its post id value here. A line will be drawn, and it
	  * will be closed once you reach the bottom of the page. */
	const lastSeenPostBeforeUpdate = ref<number | null>(null);
	/** `post.user_id` => number of posts with that ID in this thread */
	const userIdCount = ref<Record<string, number>>({});

	const replyForm = useTemplateRef('reply-form');
	const textSelection = useTextSelection();

	const autoTimer = ref<number>(10);
	const autoTimerIsEnabled = ref<boolean>(false);

	const floatingReplyBox = useTemplateRef("floating-reply-box");
	const floatingReplyBoxDragHandle = useTemplateRef("drag-handle");
	const floatingReplyBoxVisible = ref<boolean>(false);
	const { style: floatingReplyBoxStyle } = useDraggable(floatingReplyBox, {initialValue: { x: 200, y: 100 }, handle: floatingReplyBoxDragHandle, });

	/** Key is `board + postNum` concatenated */
	const peekPostCache = ref<Record<string, PostPeekBundle>>({});
	const peekPostVisible = ref<boolean>(false);
	const peekPost = ref<PostPeekBundle | undefined>(undefined);
	const peekMouseX = ref<number>(0);
	const peekMouseY = ref<number>(0);

	const userRole = ref<UserRole | undefined>(undefined);

	const galleryMode = useTemplateRef("galleryMode");

	onMounted(() => {
		const board_code: string = route.params.board_code as string;
		loadBoard(board_code);
		autoTimerCountdown();
		window.addEventListener('mousemove', updatePosition);

		userRole.value = (localStorage.getItem("role") ?? undefined) as UserRole;
	});

	onUnmounted(() => {
		window.removeEventListener('mousemove', updatePosition);
	})

	const updatePosition = (e: MouseEvent) => {
		const elemHeight: number = document.getElementById("peekElement")?.clientHeight ?? 80;
		const maxY = window.innerHeight - elemHeight;

		peekMouseX.value = e.clientX;
		peekMouseY.value = e.clientY - 64;
		if (peekMouseY.value > maxY) {
			peekMouseY.value -= elemHeight;
		}
	}

	watch(() => route.hash, (newHash) => {
		if (newHash) {
			const chunk = newHash.substring(1);

			if (chunk.startsWith("p")) {
				const highlightedPostNum = Number(chunk.substring(1));
				highlightedPost.value = highlightedPostNum
			}
		}
	});

	const setTabTitle = () => {
		let newPostCount = 0;
		if (lastSeenPostBeforeUpdate.value && posts.value) {
			for (var i = 0; i < posts.value.length; i++) {
				if (posts.value[i]?.id == lastSeenPostBeforeUpdate.value) {
					newPostCount = posts.value.length - (i + 1);
					break;
				}
			}
		}

		document.title = GetTabTitleForThread(board.value, thread.value, posts.value, newPostCount);
	}

	const loadBoard = (boardCode: string) => {
		BoardAPI.GetBoard(boardCode).then((res: AxiosResponse<ApiResponse<BoardDTO>>) => {
			board.value = res.data.data!;

			const thread_num: number = Number(route.params.thread_num);
			loadThread(boardCode, thread_num);
		}).catch((err: AxiosError) => {
			router.push("/");
		});
	}

	const loadThread = (board_code: string, thread_num: number) => {
		ThreadAPI.GetFullThreadByNum(board_code, thread_num).then((res: AxiosResponse<ApiResponse<ThreadFullDTO>>) => {
			const loadingPostsForTheFirstTime: boolean = posts.value.length == 0;

			const dto: ThreadFullDTO = res.data.data!;
			thread.value = dto.thread;
			posts.value = dto.posts.sort((a: PostDTO, b: PostDTO) => a.id - b.id);

			thread_stats.value.posts = posts.value.length;
			thread_stats.value.images = posts.value.filter((p: PostDTO) => p.filename).length;
			thread_stats.value.posters = [... new Set(posts.value.map((p: PostDTO) => p.ipv4))].length;
			thread_stats.value.page = 1;

			for (let p of posts.value) {
				processPost(p);
			}

			if (board.value!.config.ids_enabled) {
				userIdCount.value = {};
				for (let p of posts.value) {
					userIdCount.value[p.public_id!] = (userIdCount.value[p.public_id!] ?? 0) + 1;
				}
			}

			if (loadingPostsForTheFirstTime) {
				if (posts.value.length > 0) {
					lastSeenPostBeforeUpdate.value = posts.value.at(-1)!.id;
				}

				if (route.hash) {
					highlightedPost.value = Number(route.hash.substring(2));

					setTimeout(() => {
						document.getElementById(route.hash.substring(1))?.scrollIntoView();
					}, 500);
				}
			} else {

			}

			setTabTitle();
			determinePage();
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const determinePage = () => {
		ThreadAPI.GetThreadsForCatalog(board.value!.code).then((res: AxiosResponse<ApiResponse<ThreadForCatalogDTO[]>>) => {
			let threads: ThreadForCatalogDTO[] = res.data.data!;
			let pinnedThreads = (localStorage.getItem("pinned-threads") ?? "").split(",");

			SortThreadsForCatalog(threads, ThreadSortModeInCatalog.BumpOrder, board.value!.code, pinnedThreads);

			const index = threads.findIndex((t) => t.thread.id == thread.value!.id);

			if (index >= 0) {
				thread_stats.value.page = Math.floor(index / 10) + 1;
			}
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const reloadThread = () => {
		if (board.value && thread.value) {
			loadThread(board.value.code, thread.value.post_num);
		}
	}

	const onClickPostNumber = (postNum: number) => {
		if (!canReply()) {
			return;
		}

		openFloatingReplyBox();

		if (replyForm.value) {
			replyForm.value.AppendText(`>>${postNum}\n`);

			// If any selected text, add it to the reply box.
			var selectionParts: string[] = textSelection.text.value.split("\n");

			if (selectionParts.length > 1 || (selectionParts[0]?.length ?? 0 > 0)) {
				var quoted = "";
				for (let part of selectionParts) {
					quoted += ">" + part + "\n";
				}
				if (quoted.length > 0) {
					replyForm.value.AppendText(quoted);
				}
			}
		}
	}

	const onClickPostNo = (postNum: number) => {
		router.replace({'hash': `#p${postNum}`});
		highlightedPost.value = postNum;
		console.log(highlightedPost.value);
	}

	const onClickPublicId = (userId: string) => {
		highlightedPublicIDs.value[userId] = !(highlightedPublicIDs.value[userId] ?? false);
	}

	const onClickPostImage = (postId: number) => {
		imageData.value[postId]!.expanded = !imageData.value[postId]!.expanded;
	}

	const onPostLinkHover = (postLink: string) => {
		let [link_post_board, link_post_num] = SplitPostLink(postLink, board.value!.code);
		if (link_post_num == 0) return;
		// Set to true immediately. If it's set to false before GetPostPeek
		// resolves, that's fine, it should be overruled.
		peekPostVisible.value = true;
		GetPostPeek(link_post_board, link_post_num, imageData.value, peekPostCache.value).then((res: PostPeekBundle) => {
			peekPost.value = res;
		}).catch((err: any) => {
			console.error(err);
			peekPostVisible.value = false;
		});
	}

	const onPostLinkUnhover = (postLink: string) => {
		peekPost.value = undefined;
		peekPostVisible.value = false;
	}

	const onToggleSticky = (threadObj: ThreadDTO) => {
		const dto: UpdateThreadDTO = {
			sticky: !threadObj.sticky,
			locked: threadObj.locked,
			auto_cycle: threadObj.auto_cycle,
		}
		ThreadAPI.UpdateThread(threadObj.id, dto).then((res) => {
			thread.value = res.data.data!;
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const onToggleLocked = (threadObj: ThreadDTO) => {
		const dto: UpdateThreadDTO = {
			sticky: threadObj.sticky,
			locked: !threadObj.locked,
			auto_cycle: threadObj.auto_cycle,
		}
		ThreadAPI.UpdateThread(threadObj.id, dto).then((res) => {
			thread.value = res.data.data!;
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const onChangeAutoCycle = (threadObj: ThreadDTO) => {
		const dto: UpdateThreadDTO = {
			sticky: threadObj.sticky,
			locked: threadObj.locked,
			auto_cycle: threadObj.auto_cycle,
		}
		ThreadAPI.UpdateThread(threadObj.id, dto).then((res) => {
			thread.value = res.data.data!;
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const deleteThread = () => {
		if (!(thread.value && board.value)) {
			return;
		}
		ThreadAPI.DeleteThread(thread.value.id).then((_) => {
			router.push(`/${board.value!.code}`);
		}).catch((err: AxiosError) => {
			console.error(err);
		})
	}

	const archiveThread = () => {
		if (!thread.value) {
			return;
		}
		ThreadAPI.ArchiveThread(thread.value.id).then((res: AxiosResponse<ApiResponse<ThreadDTO>>) => {
			thread.value = res.data.data!;
		}).catch((err: AxiosError) => {
			console.error(err);
		})
	}

	const processPost = (post: PostDTO) => {
		ProcessPost(post, thread.value!, board.value!, imageData.value, postLinks.value, posts.value.map((p) => p.num))
		.then((res: ProcessedPost) => {
			if (res.image) {
				imageData.value[post.id] = res.image;
			}

			for (const link_post_num of res.backlinks) {
				const links = [post.num];
				backLinks.value[link_post_num] = AddRangeNoDuplicates(backLinks.value[link_post_num] ?? [], links);
			}

			postTokens.value[post.id] = res.tokens;
			for (const linkKey in res.links) {
				postLinks.value[linkKey] = res.links[linkKey]!;
			}

			// Add to peek cache. Chances are most peeked posts will be part
			// of the current thread instead of being cross-linked.
			peekPostCache.value[`${board.value!.code}/${post.num}`] = {
				post: post,
				thread: thread.value!,
				board: board.value!,
				imageData: res.image,
				tokens: res.tokens,
			};
		}).catch((err: any) => {
			console.error(err);
		});
	}

	const onPostCreated = () => {
		replyForm.value?.Clear();
		reloadThread();
	}

	const autoTimerCountdown = () => {
		setTimeout(() => {
			if (autoTimer.value <= 0) {
				autoTimer.value = 10;
				reloadThread();
			} else {
				if (autoTimerIsEnabled.value) {
					autoTimer.value -= 1;
				}
			}

			autoTimerCountdown();
		}, 1000);
	}

	const onAutoTimerToggled = (enabled: boolean) => {
		autoTimerIsEnabled.value = enabled;
	}

	const onLastPostSeenVisibilityNotify = (isVisible: boolean) => {
		if (isVisible) {
			if (posts.value.length > 0) {
				lastSeenPostBeforeUpdate.value = posts.value.at(-1)!.id;
				setTabTitle();
			}
		}
	}

	const closeFloatingReplyBox = () => {
		floatingReplyBoxVisible.value = false;
	}

	const openFloatingReplyBox = () => {
		floatingReplyBoxVisible.value = true;
	}

	const canReply = (): boolean => {
		return !thread.value!.locked && !thread.value!.archived && !board.value!.config.locked;
	}

	const openGalleryMode = () => {
		galleryMode.value!.OpenGalleryMode();
	}
</script>

<template>
	<BoardListNav :isCatalog=false />

	<template v-if="board && thread">
		<GalleryMode ref="galleryMode" :board="board" :thread="thread" :posts="posts" />
	</template>

	<!-- Peek post -->
	<template v-if="peekPostVisible && peekPost">
		<PostComponent
		class="peek"
		id="peekElement"
		:userRole="undefined"
		:style="{ transform: 'translate(' + peekMouseX + 'px,' + peekMouseY + 'px)' }"
		:board="peekPost.board"
		:thread="peekPost.thread"
		:post="peekPost.post"
		:is_highlighted="false"
		:is_op_post="false"
		:is_last_seen="false"
		:backlinks="[]"
		:image_data="undefined"
		:post_tokens="peekPost.tokens"
		:public_id_count="undefined"
		/>
	</template>

	<template v-if="board && thread">
		<BoardBanner :board="board" />

		<!-- Static reply box -->
		<div v-if="canReply()">
			<hr />
			<CreatePostForm
			v-if="!board.config.locked"
			id="static-reply-box"
			:board="board"
			:thread="thread"
			:max_size_bytes="board.config.max_file_size"
			:mime_types_allowed="board.config.mime_types_allowed"
			@postCreated="onPostCreated()"
			/>
		</div>
		<div v-else>
			<p class="red">
				Thread
				<template v-if="thread.locked">closed.</template>
				<template v-if="thread.archived">archived.</template>
				<template v-if="board.config.locked">locked with the board.</template>
				<br/>
				You my not reply at this time.
			</p>
		</div>

		<RandomBoardImageBanner />

		<!-- Floating reply box -->
		<div :style="{display: floatingReplyBoxVisible ? 'block' : 'none'}">
			<div id="floating-reply-box" ref="floating-reply-box" :style="floatingReplyBoxStyle">
				<div class="header" ref="drag-handle">
					<label>Reply to Thread No.{{thread.post_num}}</label>
					<img src="/icons/cross.png" @click="closeFloatingReplyBox" />
				</div>
				<CreatePostForm
				ref="reply-form"
				:board="board"
				:thread="thread"
				:max_size_bytes="board.config.max_file_size"
				:mime_types_allowed="board.config.mime_types_allowed"
				@postCreated="onPostCreated()"
				/>
			</div>
		</div>

		<ThreadViewNavList
		:board_code="board.code"
		jump_to_id="bottom"
		jump_to_label="Bottom"
		:thread_stats="thread_stats"
		:autoTimer="autoTimer"
		:isAutoTimerUsed="autoTimerIsEnabled"
		:sticky="thread!.sticky"
		:locked="thread!.locked"
		:showCenterElements="false"
		@updateClicked="reloadThread"
		@autoTimerToggled="onAutoTimerToggled"
		@openedGalleryMode="openGalleryMode" />

		<template v-if="thread">
			<PostComponent
			v-for="post, i of posts"
			:userRole="userRole"
			:board="board"
			:thread="thread"
			:post="post"
			:is_highlighted="highlightedPost == post.num || (highlightedPublicIDs[post.public_id ?? ''] ?? false)"
			:is_op_post="thread.post_num == post.num"
			:is_last_seen="lastSeenPostBeforeUpdate == post.id && i != posts.length - 1"
			:backlinks="backLinks[post.num] ?? []"
			:image_data="imageData[post.id]"
			:post_tokens="postTokens[post.id] ?? []"
			:public_id_count="userIdCount"
			@onClickPostNo="onClickPostNo"
			@onClickPostNumber="onClickPostNumber"
			@onClickPostImage="onClickPostImage"
			@onClickPublicId="onClickPublicId"
			@onPostLinkHover="onPostLinkHover"
			@onPostLinkUnhover="onPostLinkUnhover"
			@onToggleSticky="onToggleSticky"
			@onToggleLocked="onToggleLocked"
			@onArchive="archiveThread"
			@onDelete="deleteThread"
			@onChangeAutoCycle="onChangeAutoCycle"
			/>
		</template>

		<span v-element-visibility="onLastPostSeenVisibilityNotify"></span>
		<ThreadViewNavList
		:board_code="board.code"
		jump_to_id="top"
		jump_to_label="Top"
		:thread_stats="thread_stats"
		:autoTimer="autoTimer"
		:isAutoTimerUsed="autoTimerIsEnabled"
		:sticky="thread!.sticky"
		:locked="thread!.locked"
		:showCenterElements="true"
		@openedReplyBox="openFloatingReplyBox"
		@updateClicked="reloadThread"
		@autoTimerToggled="onAutoTimerToggled"
		@openedGalleryMode="openGalleryMode" />
	</template>

	<BoardListNav :isCatalog=false />
</template>

<style scoped>
	.error {
		color: var(--user-error-color);
	}

	#static-reply-box {
		display: block;
		text-align: center;
		width: 30%;
		margin: auto;
	}

	.red {
		text-align: center;
		font-weight: bold;
		color: var(--user-error-color);
	}

	#floating-reply-box {
		position: fixed;
		z-index: 600;
		overflow: hidden;
		background-color: var(--post-background-color);
		border: 1px solid gray;
		padding: 2px;

		.header {
			background-color: var(--background-color-accent);
			font-weight: bold;
			text-align: center;
			cursor: move;

			label {
				-webkit-user-select: none;
				-moz-user-select: none;
				-ms-user-select: none;
				user-select: none;
				cursor: move;
			}

			img {
				float: right;
				cursor: pointer;
			}
		}
	}
</style>
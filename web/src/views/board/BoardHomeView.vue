<script setup lang="ts">
	import { BoardAPI, type BoardDTO } from "@/api/board.api.ts";
	import type { ApiResponse } from '@/api/http';
	import { type PostDTO } from "@/api/post.api";
	import { ThreadAPI, type ThreadForHomeDTO, type ThreadDTO, type UpdateThreadDTO } from "@/api/thread.api.ts";
	import PostComponent from "@/components/post/PostComponent.vue";
	import { type PostImageData, type PostToken, type PostLinkToken, ProcessPost, type ProcessedPost, SplitPostLink } from "@/model/post/post.model";
	import type { AxiosError, AxiosResponse } from 'axios';
	import { onMounted, onUnmounted, reactive, ref } from 'vue';
	import { useRoute, useRouter } from "vue-router";
	import { GetPostPeek, type PostPeekBundle } from "@/model/post/post.peek";
	import { GetTabTitleForBoard } from "@/util/tab.util";
	import CreatePostForm from "@/components/post/CreatePostForm.vue";
	import type { UserRole } from "@/api/user.api";
	import BoardBanner from "@/components/board/BoardBanner.vue";
	import RandomBoardImageBanner from "@/components/board/RandomBoardImageBanner.vue";
	import BlogpostQuickList from '@/components/blogpost/BlogpostQuickList.vue';
import { Paginator } from "@/util/pagination.util";
import PaginatorComponent from "@/components/PaginatorComponent.vue";
	
	const board = ref<BoardDTO | undefined>(undefined);

	const route = useRoute();
	const router = useRouter();
	
	const threads = ref<ThreadForHomeDTO[]>([]);
	const threadsError = ref<string | undefined>(undefined);


	const getThreads = (offset: number, limit: number) => ThreadAPI.GetThreadsForHome(board.value!.code, offset, limit);
	const pagination = reactive<Paginator<ThreadForHomeDTO[]>>(new Paginator(getThreads));

	/** `post.user_id` which should be highlighted */
	const highlightedPublicIDs = ref<Record<string, boolean>>({});
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
	
	/** Key is `board + postNum` concatenated */
	const peekPostCache = ref<Record<string, PostPeekBundle>>({});
	const peekPostVisible = ref<boolean>(false);
	const peekPost = ref<PostPeekBundle | undefined>(undefined);
	const peekMouseX = ref<number>(0);
	const peekMouseY = ref<number>(0);

	const userRole = ref<UserRole | undefined>(undefined);

	onMounted(() => {
		const board_code: string = route.params.board_code as string;
		const page_num_string = route.query['page'] ?? "1";
		pagination.page = Number(page_num_string) - 1;
		loadBoard(board_code);
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
	
	const loadBoard = (boardCode: string) => {
		BoardAPI.GetBoard(boardCode).then((res: AxiosResponse<ApiResponse<BoardDTO>>) => {
			board.value = res.data.data!;
			loadThreads();
			document.title = GetTabTitleForBoard(board.value, false);
		}).catch((err: AxiosError) => {
			router.push("/");
		});
	}
	
	const loadThreads = () => {
		if (board.value == null) {
			return;
		}

		pagination.getItems()
			.then((res) => {
				threads.value = res.data.data!;

				for (let i = 0; i < threads.value.length; i++) {
					const thread: ThreadForHomeDTO = threads.value[i]!;
					for (let j = 0; j < thread.posts.length; j++) {
						processPost(thread.posts[j]!, thread);
					}
				}
			})
			.catch((_) => threadsError.value = "Could not fetch threads!");
	}

	const onClickPostNo = (post_num: number, thread: ThreadForHomeDTO) => {
		router.push(`${board.value!.code}/thread/${thread.thread.post_num}#p${post_num}`);
	}

	const onClickPostNumber = (post_num: number, thread: ThreadForHomeDTO) => {
		router.push(`${board.value!.code}/thread/${thread.thread.post_num}#p${post_num}`);
	}

	const onClickPublicId = (publicId: string) => {
		highlightedPublicIDs.value[publicId] = !(highlightedPublicIDs.value[publicId] ?? false);
	}

	const onClickPostImage = (postId: number, thread: ThreadForHomeDTO) => {
		imageData.value[postId]!.expanded = !imageData.value[postId]!.expanded;
	}

	const onToggleSticky = (threadObj: ThreadDTO) => {
		const dto: UpdateThreadDTO = {
			sticky: !threadObj.sticky,
			locked: threadObj.locked,
			auto_cycle: threadObj.auto_cycle,
		}
		ThreadAPI.UpdateThread(threadObj.id, dto).then((res) => {
			updateThreadObject(threadObj.id, res.data.data!);
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const updateThreadObject = (thread_id: number, threadObj: ThreadDTO) => {
		for (let i = 0; i < threads.value.length; i++) {
			if (threads.value[i]!.thread.id == thread_id) {
				threads.value[i]!.thread = threadObj;
				return;
			}
		}
	}

	const updatePostObject = (postDTO: PostDTO) => {
		for (let i = 0; i < threads.value.length; i++) {
			if (threads.value[i]!.thread.id == postDTO.thread_id) {
				const thread = threads.value[i]!;

				for (let j = 0; j < thread.posts.length; j++) {
					if (thread.posts[j]!.id == postDTO.id) {
						thread.posts[j] = postDTO;
						return;
					}
				}

				return;
			}
		}
	}

	const deletedPostObject = (id: number) => {
		for (let i = 0; i < threads.value.length; i++) {
			const thread = threads.value[i]!;
			for (let j = 0; j < thread.posts.length; j++) {
				if (thread.posts[j]!.id == id) {
					thread.posts.splice(j, 1);
					return;
				}
			}
		}
	}

	const onToggleLocked = (threadObj: ThreadDTO) => {
		const dto: UpdateThreadDTO = {
			sticky: threadObj.sticky,
			locked: !threadObj.locked,
			auto_cycle: threadObj.auto_cycle,
		}
		ThreadAPI.UpdateThread(threadObj.id, dto).then((res) => {
			updateThreadObject(threadObj.id, res.data.data!);
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
			updateThreadObject(threadObj.id, res.data.data!);
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const deleteThread = (thread: ThreadDTO) => {
		if (!board.value) {
			return;
		}
		ThreadAPI.DeleteThread(thread.id).then((_) => {
			loadThreads();
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const archiveThread = (thread: ThreadDTO) => {
		ThreadAPI.ArchiveThread(thread.id).then((res: AxiosResponse<ApiResponse<ThreadDTO>>) => {
			updateThreadObject(thread.id, res.data.data!);
		}).catch((err: AxiosError) => {
			console.error(err);
		});
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

	const processPost = (post: PostDTO, thread: ThreadForHomeDTO) => {
		ProcessPost(post, thread.thread, board.value!, imageData.value, postLinks.value, thread.posts.map((p) => p.num))
		.then((res: ProcessedPost) => {
			if (res.image) {
				imageData.value[post.id] = res.image;
			}

			postTokens.value[post.id] = res.tokens;

			for (const linkKey in res.links) {
				postLinks.value[linkKey] = res.links[linkKey]!;
			}

			// Add to peek cache. Chances are most peeked posts will be part
			// of the current thread instead of being cross-linked.
			peekPostCache.value[`${board.value!.code}/${post.num}`] = {
				post: post,
				thread: thread.thread,
				board: board.value!,
				imageData: res.image,
				tokens: res.tokens,
			};
		}).catch((err: any) => {
			console.error(err);
		});
	}

	const gotoPage = (page: number) => {
		router.push(`/${board.value!.code}?page=${page}`);
	}
</script>

<template>
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
		:peek="true"
		:post_tokens="peekPost.tokens"
		:public_id_count="undefined"
		/>
	</template>

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
		</div>

		<BlogpostQuickList />
		<RandomBoardImageBanner />

		<hr />
		<!-- Navigation and search -->
		[<RouterLink :to="`/${route.params.board_code}/catalog`">Catalog</RouterLink>]
		[<RouterLink :to="`/${route.params.board_code}/archive`">Archive</RouterLink>]
		[<a class="link" :href="`#bottom`">Bottom</a>]

		<hr />

		<!-- Thread list -->
		<template v-if="threadsError">
			<div class="error">{{ threadsError }}</div>
		</template>
		<template v-else>
			<div v-for="thread in threads">
				<template v-for="post, i of thread.posts">
					<PostComponent
					:userRole="userRole"
					:board="board"
					:thread="thread.thread"
					:post="post"
					:is_highlighted="(highlightedPublicIDs[post.public_id ?? ''] ?? false)"
					:is_op_post="i == 0"
					:is_last_seen="false"
					:backlinks="[]"
					:image_data="imageData[post.id]"
					:post_tokens="postTokens[post.id] ?? []"
					:public_id_count="undefined"
					:peek="false"
					@onClickPostNo="(n: number) => onClickPostNo(n, thread)"
					@onClickPostNumber="(n: number) => onClickPostNumber(n, thread)"
					@onClickPostImage="(n: number) => onClickPostImage(n, thread)"
					@onClickPublicId="onClickPublicId"
					@onPostLinkHover="onPostLinkHover"
					@onPostLinkUnhover="onPostLinkUnhover"
					@onToggleSticky="onToggleSticky"
					@onToggleLocked="onToggleLocked"
					@onArchive="archiveThread"
					@onDelete="deleteThread"
					@onChangeAutoCycle="onChangeAutoCycle"
					@onPostUpdated="(dto: PostDTO) => updatePostObject(dto)"
					@onPostDeleted="(post_id: number) => deletedPostObject(post_id)"
					/>
					<div v-if="i == 0 && thread.posts.length < thread.stats.post_count">
						<template v-if="thread.stats.image_count - thread.posts.filter((p) => p.filename).length > 0">
							{{ thread.stats.image_count - thread.posts.filter((p) => p.filename).length }} images
						</template>
						<template v-if="thread.stats.image_count - thread.posts.filter((p) => p.filename).length > 0 && thread.stats.post_count - thread.posts.length > 0">
							and
						</template>
						<template v-if="thread.stats.post_count - thread.posts.length > 0">
							{{ thread.stats.post_count - thread.posts.length }} replies
						</template>
						ommited.
						<RouterLink :to="`${board.code}/thread/${thread.thread.post_num}`">Click here</RouterLink> to view them.
					</div>
				</template>

				<hr />
			</div>
		</template>

		<!-- Pagination -->
		<span class="pagination">
			<PaginatorComponent :paginator="pagination" @gotoPage="gotoPage" emptyMessage="No threads" />
			<span class="nav">
				<!-- Navigation and search #2 -->
				[<RouterLink :to="`/${route.params.board_code}/catalog`">Catalog</RouterLink>]
				[<RouterLink :to="`/${route.params.board_code}/archive`">Archive</RouterLink>]
				[<a class="link" :href="`#top`">Top</a>]
			</span>
		</span>
	</template>
</template>

<style scoped>
	#create-thread {
		display: block;
		text-align: center;
		width: 30%;
		margin: auto;
	}

	.error {
		color: var(--user-error-color);
	}

	.currentPage {
		font-weight: bold;
	}

	.pagination {
		background-color: var(--post-background-color);
		display: inline-block;
		padding: 1em;
		border: 1px solid var(--post-border-color);

		.nav {
			margin-left: 1em;
			&* {
				margin-left: 0.2em;
			}
		}
	}
</style>
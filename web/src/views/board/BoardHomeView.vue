<script setup lang="ts">
	import { BoardAPI, type BoardDTO } from "@/api/board.api.ts";
	import type { ApiResponse } from '@/api/http';
	import { type PostDTO } from "@/api/post.api";
	import { ThreadAPI, type ThreadForHomeDTO } from "@/api/thread.api.ts";
	import PostComponent from "@/components/post/PostComponent.vue";
	import CreateThreadForm from '@/components/thread/CreateThreadForm.vue';
	import { type PostImageData, type PostToken, type PostLinkToken, ProcessPost, type ProcessedPost, SplitPostLink } from "@/model/post/post.model";
	import type { AxiosError, AxiosResponse } from 'axios';
	import { onMounted, onUnmounted, ref } from 'vue';
	import { useRoute, useRouter } from "vue-router";
	import BoardListNav from '@/components/board/BoardListNav.vue';
	import { GetPostPeek, type PostPeekBundle } from "@/model/post/post.peek";
	
	const board = ref<BoardDTO | null>(null);

	const route = useRoute();
	const router = useRouter();
	
	const threads = ref<ThreadForHomeDTO[]>([]);
	const threadsError = ref<string | undefined>(undefined);

	const page = ref<number>(0);
	const pageSize = ref<number>(2);
	const totalPages = ref<number>(0);
	const pages = ref<number[]>([]);

	/** `post.user_id` which should be highlighted */
	const highlightedUserIDs = ref<Record<string, boolean>>({});
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

	onMounted(() => {
		const board_code: string = route.params.board_code as string;
		const page_num_string = route.query['page'] ?? "1";
		page.value = Number(page_num_string) - 1;
		loadBoard(board_code);
		window.addEventListener('mousemove', updatePosition);
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
		ThreadAPI.GetThreadsForHome(board.value.code, page.value, pageSize.value).then((res : AxiosResponse<ApiResponse<ThreadForHomeDTO[]>>) => {
			if (res.data.meta) {
				totalPages.value = res.data.meta.total_pages;
				pages.value = [];
				for (let i = 1; i <= totalPages.value; i++) {
					pages.value.push(i);
				}
			}

			threads.value = res.data.data! ?? [];

			// Process posts into tokens.
			for (let i = 0; i < threads.value.length; i++) {
				const thread: ThreadForHomeDTO = threads.value[i]!;
				for (let j = 0; j < thread.posts.length; j++) {
					processPost(thread.posts[j]!, thread);
				}
			}
		}).catch((err: AxiosError) => {
			threadsError.value = "Could not fetch threads";
			console.error(err);
		});
	}

	const onClickPostNo = (post_num: number, thread: ThreadForHomeDTO) => {
		router.push(`${board.value!.code}/thread/${thread.thread.post_num}#p${post_num}`);
	}

	const onClickPostNumber = (post_num: number, thread: ThreadForHomeDTO) => {
		router.push(`${board.value!.code}/thread/${thread.thread.post_num}#p${post_num}`);
	}

	const onClickUserId = (userId: string) => {
		highlightedUserIDs.value[userId] = !(highlightedUserIDs.value[userId] ?? false);
	}

	const onClickPostImage = (postId: number, thread: ThreadForHomeDTO) => {
		imageData.value[postId]!.expanded = !imageData.value[postId]!.expanded;
	}

	const onPostLinkHover = (postLink: string) => {
		let [link_post_board, link_post_num] = SplitPostLink(postLink, board.value!.code);
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
</script>

<template>
	<BoardListNav :isCatalog=false />

	<template v-if="peekPostVisible && peekPost">
		<PostComponent
		class="peek"
		id="peekElement"
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
		:user_id_count="undefined"
		/>
	</template>

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
		@threadCreated="loadThreads()" />
		<hr />

		<!-- Navigation and search -->
		[<RouterLink :to="`/${route.params.board_code}/catalog`">Catalog</RouterLink>]
		<!-- [<RouterLink :to="`/${route.params.board_code}/archive`">Archive</RouterLink>] -->
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
					:board="board"
					:thread="thread.thread"
					:post="post"
					:is_highlighted="(highlightedUserIDs[post.user_id ?? ''] ?? false)"
					:is_op_post="i == 0"
					:is_last_seen="false"
					:backlinks="[]"
					:image_data="imageData[post.id]"
					:post_tokens="postTokens[post.id] ?? []"
					:user_id_count="undefined"
					@onClickPostNo="(n: number) => onClickPostNo(n, thread)"
					@onClickPostNumber="(n: number) => onClickPostNumber(n, thread)"
					@onClickPostImage="(n: number) => onClickPostImage(n, thread)"
					@onClickUserId="onClickUserId"
					@onPostLinkHover="onPostLinkHover"
					@onPostLinkUnhover="onPostLinkUnhover"
					/>
					<div v-if="i == 0 && thread.posts.length < thread.stats.post_count">
						{{ thread.stats.image_count - thread.posts.filter((p) => p.filename).length }} images and {{ thread.stats.post_count - thread.posts.length }} replies ommited.
						<RouterLink :to="`${board.code}/thread/${thread.thread.post_num}`">Click here</RouterLink> to view them.
					</div>
				</template>

				<hr />
			</div>
		</template>

		<!-- Pagination -->
		<div>
			<span v-if="page > 0">
				[<a :href="`/${board.code}?page=${page - 1 + 1}`">Prev</a>]
			</span>
			<span v-for="p of pages">
				[<a :href="`/${board.code}?page=${p}`" :class="{currentPage : p-1 == page}">{{p}}</a>]
			</span>
			<span v-if="page < totalPages - 1">
				[<a :href="`/${board.code}?page=${page + 1 + 1}`">Next</a>]
			</span>
		</div>

		<hr />

		<!-- Navigation and search #2 -->
		[<RouterLink :to="`/${route.params.board_code}/catalog`">Catalog</RouterLink>]
		<!-- [<RouterLink :to="`/${route.params.board_code}/archive`">Archive</RouterLink>] -->
		[<a class="link" :href="`#top`">Top</a>]
	</template>

	<BoardListNav :isCatalog=false />
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


	.error {
		color: red;
	}

	.currentPage {
		font-weight: bold;
	}

	.peek {
		position: fixed;
		z-index: 1000;
		overflow: hidden;
		pointer-events: none;
		box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
		background-color: #D6DAF0;
	}
</style>
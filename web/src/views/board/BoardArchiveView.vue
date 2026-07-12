<script setup lang="ts">
	import { BoardAPI, type BoardDTO } from '@/api/board.api';
	import type { ApiResponse } from '@/api/http';
	import { ThreadAPI, type ThreadForCatalogDTO } from '@/api/thread.api';
	import BlogpostQuickList from '@/components/blogpost/BlogpostQuickList.vue';
	import RandomBoardImageBanner from '@/components/board/RandomBoardImageBanner.vue';
	import PaginatorComponent from '@/components/PaginatorComponent.vue';
	import PostComponent from '@/components/post/PostComponent.vue';
	import { GetPostTimeReadable, SplitPostLink, type PostImageData } from '@/model/post/post.model';
	import { GetPostPeek, type PostPeekBundle } from '@/model/post/post.peek';
	import { Paginator } from '@/util/pagination.util';
	import { GetTabTitleForBoard } from '@/util/tab.util';
	import type { AxiosResponse, AxiosError } from 'axios';
	import { onMounted, onUnmounted, reactive, ref } from 'vue';
	import { useRoute, useRouter } from 'vue-router';

	const board = ref<BoardDTO | undefined>(undefined);
	const threads = ref<ThreadForCatalogDTO[]>([]);

	const route = useRoute();
	const router = useRouter();

	const getThreads = (offset: number, limit: number) => ThreadAPI.GetThreadsForArchive(board.value!.code, offset, limit);
	const pagination = reactive<Paginator<ThreadForCatalogDTO[]>>(new Paginator(getThreads));

	/** Key is `board + postNum` concatenated */
	const peekPostCache = ref<Record<string, PostPeekBundle>>({});
	const peekPostVisible = ref<boolean>(false);
	const peekPost = ref<PostPeekBundle | undefined>(undefined);
	const peekMouseX = ref<number>(0);
	const peekMouseY = ref<number>(0);
	/** `post.id` => information about the image attached to the post */
	const imageData = ref<Record<number, PostImageData>>({});

	onMounted(() => {
		const board_code: string = route.params.board_code as string;
		pagination.perPage = 4;
		const page_num_string = route.query['page'] ?? "1";
		pagination.page = Number(page_num_string);

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
			board.value = res.data.data!;
			document.title = GetTabTitleForBoard(board.value, true);
			loadThreads();
		}).catch((err: AxiosError) => {
			router.push("/");
		});
	}

	const loadThreads = () => {
		pagination.getItems()
			.then((res) => {
				threads.value = res.data.data!;
			})
			.catch((err) => console.error(err));
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

	const gotoPage = (page: number) => {
		router.push(`/${board.value!.code}/archive?page=${page}`);
		pagination.page = page;
		loadThreads();
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
		:peek="true"
		:image_data="undefined"
		:post_tokens="peekPost.tokens"
		:public_id_count="undefined"
		/>
	</template>

	<div v-if="board">
		<div id="title">
			<h1>/{{board.code}}/ - {{board.name}}
				<img v-if="board.config.locked" src="/icons/lock.png" title="Board locked for further posts" class="icon" />
			</h1>
			<small>{{board.description}}</small>
		</div>

		<BlogpostQuickList />
		<RandomBoardImageBanner />

		<hr />

		<!-- Navigation and search -->
		[<RouterLink :to="`/${route.params.board_code}`">Return</RouterLink>]
		[<RouterLink :to="`/${route.params.board_code}/catalog`">Catalog</RouterLink>]
		[<a class="link" :href="`#bottom`">Bottom</a>]
		<hr />

		<div class="title">
			<h2>Archive</h2>
		</div>

		<table>
			<tbody>
				<tr>
					<th class="w10">No.</th>
					<th class="w40">Excerpt</th>
					<th>Posts</th>
					<th>Images</th>
					<th>Users</th>
					<th class="w10">Created At</th>
					<th class="w10">Archived At</th>
				</tr>
				<tr v-for="thread of threads">
					<td>
						<RouterLink
						:to="`/${board.code}/thread/${thread.post.num}`"
						@pointerenter="onPostLinkHover(`>>${thread.post.num}`)"
						@pointerleave="onPostLinkUnhover(`>>${thread.post.num}`)">
						&gt;&gt;{{ thread.post.num }}
					</RouterLink>
				</td>
				<td>{{ thread.post.content }}</td>
				<td>{{ thread.stats.post_count }}</td>
				<td>{{ thread.stats.image_count }}</td>
				<td>{{ thread.stats.user_count }}</td>
				<td>{{ GetPostTimeReadable(thread.post.created_at) }}</td>
				<td>{{ GetPostTimeReadable(thread.stats.last_bump) }}</td>
			</tr>
		</tbody>
	</table>

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
</div>
</template>

<style scoped>
	.title {
		text-align: center;
		h1, h2 {
			color: var(--banner-title-color);
		}
	}

	table {
		width: 80%;
		margin: auto;

		tr:nth-child(even) {
			background-color: var(--background-color-darker);
		}

		tr:nth-child(odd) {
			background-color: var(--background-color);
		}

		th {
			border: 1px solid black;
			background-color: var(--background-color-accent);
			font-weight: bold;
			height: 2em;
		}

		td {
			padding: 0.2em;
		}
	}

	.w10 { width: 10%; }
	.w40 { width: 40%; }
</style>
<script setup lang="ts">
	import { ref, onMounted, watch, useTemplateRef } from 'vue';
	import { useRoute, useRouter } from "vue-router";
	import { BoardAPI, type BoardDTO } from "@/api/board.api.ts";
	import { ThreadAPI, type ThreadDTO, type ThreadFullDTO } from "@/api/thread.api.ts";
	import { type PostDTO } from "@/api/post.api.ts";
	import ThreadViewNavList from "@/components/thread/ThreadViewNavList.vue";
	import { type ThreadStats } from "@/model/thread/thread.model.ts";
	import CreatePostForm from '@/components/post/CreatePostForm.vue';
	import type { AxiosError, AxiosResponse } from 'axios';
	import type { ApiResponse } from '@/api/http';
	import { type PostImageData, type PostToken, type PostLinkToken, type ProcessedPost, ProcessPost } from '@/model/post/post.model';
	import { vElementVisibility } from '@vueuse/components';
	import PostComponent from '@/components/post/PostComponent.vue';
	import { AddRangeNoDuplicates } from '@/util/various.util';
	import BoardListNav from '@/components/board/BoardListNav.vue';

	const route = useRoute();
	const router = useRouter();

	const board = ref<BoardDTO | null>(null);
	const thread = ref<ThreadDTO | null>(null);
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
	const highlightedUserIDs = ref<Record<string, boolean>>({});
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

	const autoTimer = ref<number>(10);
	const autoTimerIsEnabled = ref<boolean>(false);

	onMounted(() => {
		const board_code: string = route.params.board_code as string;
		loadBoard(board_code);
		autoTimerCountdown();
	});

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
		const websiteTitle = "ImmyChan";
		let title: string = "";

		if (lastSeenPostBeforeUpdate.value && posts.value) {
			// How many new posts since last seen post
			let newPostCount = 0;
			for (var i = 0; i < posts.value.length; i++) {
				if (posts.value[i]?.id == lastSeenPostBeforeUpdate.value) {
					newPostCount = posts.value.length - (i + 1);
					break;
				}
			}

			if (newPostCount > 0) {
				title += `(${newPostCount}) `;
			}
		}

		if (board.value) {
			title += `/${board.value?.code}/ - `;
		}
		if (thread.value && posts.value) {
			if (thread.value?.subject) {
				title += thread.value.subject;
				title += " - ";
			} else if (posts.value[0]?.content) {
				title += posts.value[0].content;
				title += " - "
			}
		}
		if (board.value) {
			title += `${board.value.name} - `;
		}

		title += websiteTitle;

		document.title = title;
	}

	const loadBoard = (boardCode: string) => {
		BoardAPI.GetBoard(boardCode).then((res: AxiosResponse<ApiResponse<BoardDTO>>) => {
			board.value = res.data.data;

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
					userIdCount.value[p.user_id!] = (userIdCount.value[p.user_id!] ?? 0) + 1;
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
		if (replyForm.value) {
			replyForm.value.AppendText(`>>${postNum}\n`);
		}
	}

	const onClickPostNo = (postNum: number) => {
		router.replace({'hash': `#p${postNum}`});
		highlightedPost.value = postNum;
	}

	const onClickUserId = (userId: string) => {
		highlightedUserIDs.value[userId] = !(highlightedUserIDs.value[userId] ?? false);
	}

	const onClickPostImage = (postId: number) => {
		imageData.value[postId]!.expanded = !imageData.value[postId]!.expanded;
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
</script>

<template>
	<BoardListNav :isCatalog=false />

	<template v-if="board && thread">
		<CreatePostForm
		ref="reply-form"
		:thread_id="thread.id"
		:max_size_bytes="board.config.max_file_size"
		:mime_types_allowed="board.config.mime_types_allowed"
		@postCreated="onPostCreated()"
		/>

		<ThreadViewNavList
		:board_code="board.code"
		jump_to_id="bottom"
		jump_to_label="Bottom"
		:thread_stats="thread_stats"
		:autoTimer="autoTimer"
		:isAutoTimerUsed="autoTimerIsEnabled"
		@updateClicked="reloadThread"
		@autoTimerToggled="onAutoTimerToggled" />

		<template v-if="thread">
			<PostComponent
			v-for="post, i of posts"
			:board="board"
			:thread="thread"
			:post="post"
			:is_highlighted="highlightedPost == post.num || (highlightedUserIDs[post.user_id ?? ''] ?? false)"
			:is_op_post="thread.post_num == post.num"
			:is_last_seen="lastSeenPostBeforeUpdate == post.id && i != posts.length - 1"
			:backlinks="backLinks[post.num] ?? []"
			:image_data="imageData[post.id]"
			:post_tokens="postTokens[post.id] ?? []"
			:post_links="postLinks"
			:user_id_count="userIdCount"
			@onClickPostNo="onClickPostNo"
			@onClickPostNumber="onClickPostNumber"
			@onClickPostImage="onClickPostImage"
			@onClickUserId="onClickUserId"
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
		@updateClicked="reloadThread"
		@autoTimerToggled="onAutoTimerToggled" />
	</template>

	<BoardListNav :isCatalog=false />
</template>

<style scoped>
	.error {
		color: red;
	}
</style>
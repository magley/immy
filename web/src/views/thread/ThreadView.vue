<script setup lang="ts">
	import { ref, onMounted, watch, useTemplateRef } from 'vue';
	import { useRoute, useRouter } from "vue-router";
	import { BoardAPI, type BoardDTO } from "@/api/board.api.ts";
	import { ThreadAPI, type ThreadDTO, type ThreadFullDTO } from "@/api/thread.api.ts";
	import { PostAPI, type PostDTO } from "@/api/post.api.ts";
	import ThreadViewNavList from "@/components/thread/ThreadViewNavList.vue";
	import { type ThreadStats } from "@/model/thread/thread.model.ts";
	import CreatePostForm from '@/components/post/CreatePostForm.vue';
	import { CdnAPI } from '@/api/cdn.api';
	import { GetFileSizeByteString } from '@/util/file.util';
	import type { AxiosError, AxiosResponse } from 'axios';
	import type { ApiResponse } from '@/api/http';
	import { GetPostTimeReadable, type PostImageData, ParsePostTokens, type PostToken, type PostLinkToken } from '@/model/post/post.model';

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
		highlightedPost.value = postNum;
	}

	const onClickPostImage = (postId: number, post: PostDTO) => {
		imageData.value[postId]!.expanded = !imageData.value[postId]!.expanded;
	}

	const processPost = (post: PostDTO) => {
		if (post.filename && !imageData.value[post.id]) {
			// Create an ImageData object for each image.
			const img = new Image();
			img.src = CdnAPI.GetPostImageURI(post)!;
			img.onload = () => {
				imageData.value[post.id] = {
					postId: post.id,
					expanded: false,
					width: img.naturalWidth,
					height: img.naturalHeight,
				};
			}
		}

		postTokens.value[post.id] = ParsePostTokens(post.content);
		for (let tok of postTokens.value[post.id]!) {
			if (tok.kind == 'link') {
				// Before the proper routes are attributed to each link, add a
				// dummy '#' href for each of the links.
				tok.href = '#';
			}
		}

		const boardCode: string = route.params.board_code as string;

		for (let tok of postTokens.value[post.id]!) {
			if (tok.kind == 'link') {
				if (tok.text in postLinks.value && postLinks.value[tok.text]!.href != '#') {
					// Copy relevant fields from the reference token that's cached in the `postLinks` dict.
					const refToken: PostLinkToken = postLinks.value[tok.text]!;
					tok.href = refToken.href;
					tok.local = refToken.local;
					tok.fail = refToken.fail;
					continue;
				}

				// Split into `link_post_board` and `link_post_num`.

				let link_post_board = boardCode;
				let link_post_num = 0;
				const link_text = tok.text.substring(2);

				if (link_text[0] == '/') {
					const j = link_text.indexOf('/', 1);

					if (j > 0) {
						link_post_board = link_text.substring(1, j);
						link_post_num = Number(link_text.substring(j + 1));
					}
				} else {
					link_post_num = Number(link_text);
				}

				// Check if the link points to a post in this thread.

				let post_is_local = false;
				if (link_post_board == boardCode) {
					for (let p of posts.value) {
						if (p.num == link_post_num) {
							post_is_local = true;
							break;
						}
					}
				}
				
				if (post_is_local) {
					tok.local = true;
					tok.href = `#p${link_post_num}`;

					// Add backlink.
					if (!backLinks.value[link_post_num]) {
						backLinks.value[link_post_num] = [];
					}
					if (!backLinks.value[link_post_num]!.includes(post.num)) {
						backLinks.value[link_post_num]!.push(post.num);
					}
				} else {
					tok.local = false;
					
					// It's in another thread, so fetch which thread it is.
					PostAPI.GetPostByNum(link_post_board, link_post_num).then((res: AxiosResponse<ApiResponse<PostDTO>>) => {
						const post: PostDTO = res.data.data!;
						tok.href = `/${link_post_board}/thread/${post.thread_num}#p${link_post_num}`;
					}).catch((err: AxiosError) => {
						tok.fail = true;
						console.error(err);
					});		
				}	

				// Cache the link token.
				postLinks.value[tok.text] = tok as PostLinkToken;
			}
		}
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
</script>

<template>
	<template v-if="board && thread">
		<CreatePostForm ref="reply-form" :thread_id="thread.id" :max_size_bytes="1*1024*1024" @postCreated="onPostCreated()"></CreatePostForm>
		
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
			<div :id="`p${post.num}`" class="postContainer" v-for="post, i of posts">
				<span v-if="thread.post_num != post.num" class="sideArrows"> &gt;&gt; </span>
				<span class="post" :class="{ highlightedPost: highlightedPost == post.num, opPost: thread.post_num == post.num }">
					<div class="post-header">
						<span class="subject" v-if="thread.subject && thread.post_num == post.num">{{ thread.subject }}</span>
						<span class="username">{{ post.name ? post.name : "Anonymous" }}</span>
						<span class="tripcode" v-if="post.tripcode">{{ post.tripcode }}</span>
						<span class="date">{{ GetPostTimeReadable(post.created_at) }}</span>
						<span class="postno"><a @click.prevent="onClickPostNo(post.num)" href="#" class="postNumLink">No.</a></span>
						<span class="postnum"><a @click.prevent="onClickPostNumber(post.num)" href="#" class="postNumLink">{{ post.num }}</a></span>
						<span class="dropdown">&#9654;</span>
						<span class="backlink-container" v-if="backLinks[post.num]">
							<a :href="`#p${num}`" class="backlink" v-for="num of backLinks[post.num]">&gt;&gt;{{num}}</a>
						</span>
					</div>
					
					<div class="post-body-container">
						<span v-if="post.filename" class="post-file-container">
							<div class="post-file-container-header">
								File: <a :href="CdnAPI.GetPostImageURI(post)" target="_blank">{{ post.src_filename }}</a>
								({{GetFileSizeByteString(post.filesize)}}, {{imageData[post.id]?.width}}x{{imageData[post.id]?.height}})
							</div>

							<a :href="CdnAPI.GetPostImageURI(post)" target="_blank" @click.prevent class="post-file-link">
								<!-- Thumbnail or real image. -->
								<img v-if="imageData[post.id]?.expanded"
								:src="CdnAPI.GetPostImageURI(post)"
								@click="onClickPostImage(post.id, post)"
								class="post-image-full" />
								<img v-else
								:src="CdnAPI.GetPostImageThumbnailURI(post)"
								@click="onClickPostImage(post.id, post)"
								class="post-image-thumb" />
							</a>
						</span>
						<div v-else class="post-no-file">
						</div>

						<span class="post-body">
							<span v-for="token of postTokens[post.id]">
								<template v-if="token.kind == 'text'">
									{{token.text}}
								</template>
								<template v-else-if="token.kind == 'link'">
									<template v-if="token.local">
										<a :href="`${postLinks[token.text]!.href}`" :class="{strikethrough: token.fail}" class="postRef">
											{{token.text}}
										</a>
									</template>
									<template v-else>
										<RouterLink :to="`${postLinks[token.text]!.href}`">
											<span :class="{strikethrough: token.fail}" class="postRef">
												{{token.text}} →
											</span>
										</RouterLink>
									</template>
								</template>
							</span>
						</span>
					</div>
				</span>
			</div>
		</template>
		
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
</template>

<style scoped>	
	.postContainer {
		display: flex;
		gap: 2px;
		margin-top: 0.2em;
		.sideArrows {
			
		}
		
		.post {
			background-color: #D6DAF0;
			padding-top: 0.25em;
			padding-bottom: 1em;
			padding-left: 1em;
			padding-right: 1em;

			&.opPost {
				background-color: #EEF2FF !important;
			}
			
			.post-header {
				span {
					margin-right: 0.25em;
				}
				
				.backlink-container {
					.backlink {
						font-size: small;
						margin-right: 0.25em;
					}
				}		
			}
			
			.post-body-container {
				.post-body {
					display: inline-block;
					margin-left: 1em;
					white-space: pre-wrap;
					word-wrap: break-word;
				}

				.post-no-file {
					margin-top: 0.5em;
				}

				.post-file-container {
					.post-file-container-header {
						display: block;
						margin-bottom: 0.5em;
					}

					.post-file-link {
						img {
							cursor: pointer;

							&.post-image-full {
								display: block;
								max-height: 100%;
								max-width: 100%;
							}

							&.post-image-thumb {
								display: inline;
								max-width: 40%;
								max-height: 40%;
								vertical-align: top;
							}
						}
					}
				}
			}
		}
		
		.highlightedPost {
			background-color: #d6bad0;
		}
	}
	
	.username, .tripcode {
		color: #157743;
		font-weight: bolder;
	}
	
	.subject {
		color: #0F0C5D;
		font-weight: bolder;
	}
	
	a {
		color: #34345c;
	}
	
	a:hover {
		color: #DD0000;
	}

	.postNumLink {
		color: black;
		text-decoration: none;
	}
	
	.postNumLink:hover {
		color: #DD0000;
	}
	
	.strikethrough {
		text-decoration: line-through;
	}
	
	.error {
		color: red;
	}
</style>
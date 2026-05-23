<script setup lang="ts">
	import type { BoardDTO } from '@/api/board.api';
	import { CdnAPI } from '@/api/cdn.api';
	import type { PostDTO } from '@/api/post.api';
	import type { ThreadDTO } from '@/api/thread.api';
	import { GetPostTimeReadable, type PostImageData, type PostLinkToken, type PostToken } from '@/model/post/post.model';
	import { GetFileSizeByteString } from '@/util/file.util';

	interface PostComponentProps {
		board: BoardDTO;
		thread: ThreadDTO;
		post: PostDTO;

		is_highlighted: boolean;
		is_op_post: boolean;
		is_last_seen: boolean;

		backlinks: number[];
		image_data: PostImageData | undefined;
		post_tokens: PostToken[];

		post_links: Record<string, PostLinkToken>;
	}

	const props = defineProps<PostComponentProps>()
	const emit = defineEmits(['onClickPostNo', 'onClickPostNumber', 'onClickPostImage']);

	const onClickPostNo = (post_num: number) => {
		emit('onClickPostNo', post_num);
	}
	const onClickPostNumber = (post_num: number) => {
		emit('onClickPostNumber', post_num);
	}

	const onClickPostImage = (post_id: number) => {
		emit('onClickPostImage', post_id);
	}
</script>

<template>
	<div :id="`p${post.num}`" class="postContainer">
		<span v-if="thread.post_num != post.num" class="sideArrows"> &gt;&gt; </span>
		<span class="post" :class="{
			highlightedPost: is_highlighted,
			opPost: is_op_post,
			lastSeenPost: is_last_seen,
		}">

		<div class="post-header">
			<span class="subject" v-if="thread.subject && thread.post_num == post.num">{{ thread.subject }}</span>
			<span class="username">{{ post.name ? post.name : "Anonymous" }}</span>
			<span class="tripcode" v-if="post.tripcode">{{ post.tripcode }}</span>
			<span class="date">{{ GetPostTimeReadable(post.created_at) }}</span>
			<span class="postno"><a @click.prevent="onClickPostNo(post.num)" href="#" class="postNumLink">No.</a></span>
			<span class="postnum"><a @click.prevent="onClickPostNumber(post.num)" href="#" class="postNumLink">{{ post.num }}</a></span>
			<span class="dropdown">&#9654;</span>
			<span class="backlink-container" v-if="backlinks">
				<a :href="`#p${num}`" class="backlink" v-for="num of backlinks">&gt;&gt;{{num}}</a>
			</span>
		</div>

		<div class="post-body-container">
			<span v-if="post.filename" class="post-file-container">
				<div class="post-file-container-header">
					File: <a :href="CdnAPI.GetPostImageURI(post)" target="_blank">{{ post.src_filename }}</a>
					({{GetFileSizeByteString(post.filesize)}}, {{image_data?.width}}x{{image_data?.height}})
				</div>

				<a :href="CdnAPI.GetPostImageURI(post)" target="_blank" @click.prevent class="post-file-link">
					<!-- Thumbnail or real image. -->
					<img v-if="image_data?.expanded"
					:src="CdnAPI.GetPostImageURI(post)"
					@click="onClickPostImage(post.id)"
					class="post-image-full" />
					<img v-else
					:src="CdnAPI.GetPostImageThumbnailURI(post)"
					@click="onClickPostImage(post.id)"
					class="post-image-thumb" />
				</a>
			</span>
			<div v-else class="post-no-file">
			</div>

			<span class="post-body">
				<span v-for="token of post_tokens">
					<template v-if="token.kind == 'text'">
						{{token.text}}
					</template>
					<template v-else-if="token.kind == 'link'">
						<template v-if="token.local">
							<a :href="`${post_links[token.text]!.href}`" :class="{strikethrough: token.fail}" class="postRef">
								{{token.text}}
							</a>
						</template>
						<template v-else>
							<RouterLink :to="`${post_links[token.text]!.href}`">
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

<style scoped>
	.postContainer {
		display: flex;
		gap: 2px;
		margin-top: 0.2em;
/*
		.sideArrows {

		}
*/
		.lastSeenPost {
			border-bottom: 2px solid red;
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

				.username, .tripcode {
					color: #157743;
					font-weight: bolder;
				}

				.subject {
					color: #0F0C5D;
					font-weight: bolder;
				}

			}

			.post-body-container {
				.post-body {
					display: inline-block;
					margin-left: 1em;
					white-space: pre-wrap;
					word-wrap: break-word;

					.postNumLink {
						color: black;
						text-decoration: none;
					}

					.postNumLink:hover {
						color: #DD0000;
					}
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

	.error {
		color: red;
	}

	.strikethrough {
		text-decoration: line-through;
	}
</style>
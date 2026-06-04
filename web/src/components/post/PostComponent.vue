<script setup lang="ts">
	import type { BoardDTO } from '@/api/board.api';
	import { CdnAPI } from '@/api/cdn.api';
	import type { PostDTO } from '@/api/post.api';
	import type { ThreadDTO } from '@/api/thread.api';
	import { GetPostTimeReadable, type PostImageData, type PostLinkToken, type PostToken } from '@/model/post/post.model';
	import { GetFileSizeByteString, GetMimeTypeFromFilename } from '@/util/file.util';
	import { GetUserIdColorBackground, GetUserIdColorForeground } from '@/util/various.util';

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

		/** Dictionary that maps user IDs to the number of posts made by that
		  * ID in the current thread. Optional, if it's not specified then
		  * this info won't be available inside the post. */
		user_id_count: Record<string, number> | undefined;
	}

	const props = defineProps<PostComponentProps>()
	const emit = defineEmits(['onClickPostNo', 'onClickPostNumber', 'onClickPostImage', 'onClickUserId', 'onPostLinkHover', 'onPostLinkUnhover']);

	const onClickPostNo = (post_num: number) => {
		emit('onClickPostNo', post_num);
	}
	const onClickPostNumber = (post_num: number) => {
		emit('onClickPostNumber', post_num);
	}
	const onClickPostImage = (post_id: number) => {
		emit('onClickPostImage', post_id);
	}
	const onClickUserId = (user_id: string) => {
		emit('onClickUserId', user_id);
	}
	const onPostLinkHover = (postLink: any) => {
		emit('onPostLinkHover', postLink);
	}
	const onPostLinkUnhover = (postLink: any) => {
		emit('onPostLinkUnhover', postLink);
	}

	const isPostImage = (post: PostDTO) => {
		return GetMimeTypeFromFilename(post.filename).startsWith('image/');
	}
	const isPostVideo = (post: PostDTO) => {
		return GetMimeTypeFromFilename(post.filename).startsWith('video/');
	}

	const getTrailingWhitespace = (s: string): string => {
		const trimmed: string = s.trim();
		const trimmedStart: number = s.indexOf(trimmed);
		return s.substring(trimmedStart + trimmed.length);
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
			<span
				v-if="board.config.ids_enabled && post.user_id"
				class="userid"
				:style="{backgroundColor: GetUserIdColorBackground(post.user_id), color: GetUserIdColorForeground(post.user_id)}"
				@click="onClickUserId(post.user_id)">
				ID:{{post.user_id}}
				<template v-if="user_id_count"> ({{user_id_count[post.user_id]}}) </template>
			</span>
			<span class="date">{{ GetPostTimeReadable(post.created_at) }}</span>
			<span class="postno"><a @click.prevent="onClickPostNo(post.num)" href="#" class="postNumLink">No.</a></span>
			<span class="postnum"><a @click.prevent="onClickPostNumber(post.num)" href="#" class="postNumLink">{{ post.num }}</a></span>
			<span class="dropdown">&#9654;</span>
			<span class="backlink-container" v-if="backlinks">
				<a
					:href="`#p${num}`"
					class="backlink" v-for="num of backlinks"
					@pointerenter="onPostLinkHover(`>>${num}`)"
					@pointerleave="onPostLinkUnhover(`>>${num}`)"
				>&gt;&gt;{{num}}</a>
			</span>
		</div>

		<div class="post-body-container">
			<span v-if="post.filename" class="post-file-container">
				<div class="post-file-container-header">
					File: <a :href="CdnAPI.GetPostImageURI(post)" target="_blank">{{ post.src_filename }}</a>
					({{GetFileSizeByteString(post.filesize)}}, {{post.img_width}}x{{post.img_height}})
					<template v-if="image_data?.expanded && isPostVideo(post)">
						- <a href="#" @click.prevent="onClickPostImage(post.id)">[Close]</a>
					</template>
				</div>

				<a :href="CdnAPI.GetPostImageURI(post)" target="_blank" @click.prevent class="post-file-link">
					<!-- Thumbnail or real image. -->
					<img v-if="!image_data?.expanded"
					:src="CdnAPI.GetPostImageThumbnailURI(post)"
					@click="onClickPostImage(post.id)"
					class="post-image-thumb" />
					<template v-else>
						<img v-if="isPostImage(post)"
							:src="CdnAPI.GetPostImageURI(post)"
							@click="onClickPostImage(post.id)"
							class="post-image-full" />
						<video v-if="isPostVideo(post)"
							:width="post.img_width"
							:height="post.img_height"
							controls
							autoplay
							class="post-image-full">
							<source :src="CdnAPI.GetPostImageURI(post)" :type="GetMimeTypeFromFilename(post.filename)" />
							Your browser does not support the video tag.
						</video>
					</template>
				</a>
			</span>
			<div v-else class="post-no-file">
			</div>

			<span class="post-body">
				<span v-for="token of post_tokens">
					<span v-if="token.kind == 'text'"
					:class="{'redtext': token.type=='redtext', 'greentext': token.type=='greentext'}">
						{{token.text}}
					</span>
					<span v-else-if="token.kind == 'link'">
						<a :href="`${token.href}`" :class="{strikethrough: token.fail}" class="postRef" @pointerenter="onPostLinkHover(token.text)" @pointerleave="onPostLinkUnhover(token.text)">
							{{token.text.trim()}}<template v-if="!token.local"> →</template>
						</a>
						{{getTrailingWhitespace(token.text)}}
					</span>
					<span v-else-if="token.kind == 'semantic'">
						<template v-if="token.type == 'math'">
							<vue-latex :expression="token.text" />
						</template>
						<template v-else-if="token.type == 'code'">
							<highlightjs autodetect :code="token.text" />
						</template>
					</span>
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
			border-bottom: 2px solid var(--unread-posts-separator);
		}

		.post {
			background-color: var(--post-background-color);
			padding-top: 0.25em;
			padding-bottom: 1em;
			padding-left: 1em;
			padding-right: 1em;

			&.opPost {
				background-color: var(--background-color) !important;
			}

			.post-header {
				span {
					margin-right: 0.25em;
				}

				.backlink-container {
					.backlink {
						font-size: small;
						margin-right: 0.25em;
						color: var(--link-color) !important;
					}
					.backlink:hover {
						color: var(--link-hover-color) !important;

					}
				}

				.postno a, .postnum a {
					color: black;
					text-decoration: none !important;
				}

				.postno a:hover, .postnum a:hover {
					color: var(--link-hover-color);
				}

				.username, .tripcode {
					color: var(--username-color);
					font-weight: bolder;
				}

				.userid {
					font-size: smaller;
					padding: 2px;
					border-radius: 10px;
					cursor: pointer;
				}

				.subject {
					color: var(--thread-subject-color);
					font-weight: bolder;
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
						cursor: default;
						img {
							cursor: pointer;

							&.post-image-thumb {
								display: inline;
								max-width: 40%;
								max-height: 40%;
								vertical-align: top;
							}

							&.post-image-full {
								display: block;
								max-height: 100%;
								max-width: 100%;
							}
						}

						video {
							&.post-image-full {
								display: block;
								max-height: 100%;
								max-width: 100%;
							}
						}
					}
				}
			}
		}

		.highlightedPost {
			background-color: var(--highlighted-post-background-color);
		}
	}

	.error {
		color: var(--user-error-color);
	}

	.strikethrough {
		text-decoration: line-through;
	}
</style>
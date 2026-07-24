<script setup lang="ts">
	import type { BoardDTO } from '@/api/board.api';
	import { CdnAPI } from '@/api/cdn.api';
	import { PostAPI, type PostDTO, type UpdatePostDTO } from '@/api/post.api';
	import type { ThreadDTO } from '@/api/thread.api';
	import { UserRole } from '@/api/user.api';
	import { GetPostTimeReadable, type PostImageData, type PostLinkToken, type PostToken } from '@/model/post/post.model';
	import { GetFileSizeByteString, GetMimeTypeFromFilename } from '@/util/file.util';
	import { GetPublicIdColorBackground, GetPublicIdColorForeground } from '@/util/various.util';
	import { onClickOutside } from '@vueuse/core';
	import type { AxiosError } from 'axios';
	import { ref, useTemplateRef, type ShallowRef, useId, nextTick, onMounted } from 'vue';
	import CreateBanComponent from '../ban/CreateBanComponent.vue';
	import { type Filter, FilterAction, GetFilterMatchingPost, LoadFilters } from '@/model/filter/filter.model.ts';
import { AppEvents, EventBus } from '@/util/eventBus.util.ts';

	const id = useId();
	const filterApplied = ref<Filter | null>(null);
	interface PostComponentProps {
		userRole: UserRole | undefined;

		peek: boolean | undefined;

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
		public_id_count: Record<string, number> | undefined;
	}

	const props = defineProps<PostComponentProps>()
	const emit = defineEmits([
		'onClickPostNo',
		'onClickPostNumber',
		'onClickPostImage',
		'onClickPublicId',
		'onPostLinkHover',
		'onPostLinkUnhover',
		'onToggleSticky',
		'onToggleLocked',
		'onArchive',
		'onDelete',
		'onChangeAutoCycle',
		'onPostUpdated',
		'onPostDeleted',
	]);

	const onClickPostNo = (post_num: number) => {
		emit('onClickPostNo', post_num);
	}
	const onClickPostNumber = (post_num: number) => {
		emit('onClickPostNumber', post_num);
	}
	const onClickPostImage = (post_id: number) => {
		emit('onClickPostImage', post_id);
	}
	const onClickPublicId = (user_id: string) => {
		emit('onClickPublicId', user_id);
	}
	const onPostLinkHover = (postLink: any) => {
		emit('onPostLinkHover', postLink);
	}
	const onPostLinkUnhover = (postLink: any) => {
		emit('onPostLinkUnhover', postLink);
	}
	const onToggleSticky = (thread: ThreadDTO) => {
		emit('onToggleSticky', thread);
	}
	const onToggleLocked = (thread: ThreadDTO) => {
		emit('onToggleLocked', thread);
	}
	const onArchive = (thread: ThreadDTO) => {
		emit('onArchive', thread);
	}
	const onDelete = (thread: ThreadDTO) => {
		emit('onDelete', thread);
	}
	const onChangeAutoCycle = (thread: ThreadDTO) => {
		emit('onChangeAutoCycle', thread);
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

	const canStickyThread = () => {
		return props.userRole == UserRole.Moderator || props.userRole == UserRole.Admin;
	}
	const canLockThread = () => {
		return props.userRole == UserRole.Moderator || props.userRole == UserRole.Admin;
	}
	const canArchiveThread = () => {
		return props.userRole == UserRole.Moderator || props.userRole == UserRole.Admin || props.userRole == UserRole.Janitor;
	}
	const canDeleteThread = () => {
		return props.userRole == UserRole.Moderator || props.userRole == UserRole.Admin || props.userRole == UserRole.Janitor;
	}
	const canChangeThreadAutoCycle = () => {
		return props.userRole == UserRole.Moderator || props.userRole == UserRole.Admin;
	}

	const clickBan = async (post: PostDTO) => {
		popupBan.value.visible =  !popupBan.value.visible;
		await nextTick();
		popupBan.value.set_position();
	}

	const clickDelete = async (post: PostDTO) => {
		popupDelete.value.visible =  !popupDelete.value.visible;
		await nextTick();
		popupDelete.value.set_position();
	}

	const clickEdit = (post: PostDTO) => {
		if (!isEditingPost.value) { beginEditingPost(); }
		else { cancelEditingPost(); }
	}

	class Popup {
		visible: boolean = false;
		elemRef: ShallowRef<HTMLDivElement | null>;
		elemPosID: string;

		constructor(ref: ShallowRef<HTMLDivElement | null>, elemPosID: string) {
			this.elemRef = ref;
			this.elemPosID = elemPosID;
		}

		set_position = (offsetX: number = 16, offsetY: number = 24) => {
			if (this.elemRef.value) {
				const elemPos = document.getElementById(this.elemPosID)!;
				const rect = elemPos.getBoundingClientRect();

				this.elemRef.value.style.left = `${rect.left + window.scrollX + offsetX}px`;
				this.elemRef.value.style.top = `${rect.top + window.scrollY + offsetY}px`;
			}
		}
	}

	onMounted(() => {
		refreshFilter();

		EventBus.on(AppEvents.FiltersRefreshed, refreshFilter);
	});

	const popupDeleteRef = useTemplateRef("delete-popup");
	const popupDelete = ref<Popup>(new Popup(popupDeleteRef, `delete-post-btn-${id}`));
	onClickOutside(popupDeleteRef, event => {
		if (!popupDelete.value.visible) { return; }
		setTimeout(() => { popupDelete.value.visible = false; }, 10);
	});

	const deletePostFile = () => {
		const dto: UpdatePostDTO = {
			name: null,
			tripcode: null,
			sage: null,
			content: null,
			filename: "",
			html: null
		}

		PostAPI.UpdatePost(props.post.id, dto).then((res) => {
			onPostUpdated(res.data.data!);
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const deletePost = () => {
		PostAPI.DeletePost(props.post.id).then((res) => {
			onPostDeleted(props.post.id);
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const onPostUpdated = (postDTO: PostDTO) => {
		emit("onPostUpdated", postDTO);
		popupDelete.value.visible = false;
	}
	const onPostDeleted = (id: number) => {
		emit("onPostDeleted", id);
		popupDelete.value.visible = false;
	}

	const isEditingPost = ref<boolean>(false);
	const editingPostContent = ref<string>("");
	const editingPostHtml = ref<string>("");

	const beginEditingPost = () => {
		isEditingPost.value = true;
		editingPostContent.value = props.post.content;
		editingPostHtml.value = props.post.html;
	}
	const cancelEditingPost = () => {
		isEditingPost.value = false;
	}
	const submitEditingPost = () => {
		let dto: UpdatePostDTO = {
			name: null,
			tripcode: null,
			sage: null,
			content: editingPostContent.value,
			filename: null,
			html: editingPostHtml.value
		};
		if (editingPostContent.value == props.post.content) { dto.content = null; }
		if (editingPostHtml.value == props.post.html) { dto.html = null; }

		PostAPI.UpdatePost(props.post.id, dto).then((res) => {
			onPostUpdated(res.data.data!);
			isEditingPost.value = false;
		}).catch((err: AxiosError) => {
			console.error(err);
		})
	}

	const popupBanRef = useTemplateRef("ban-popup");
	const popupBan = ref<Popup>(new Popup(popupBanRef, `ban-post-btn-${id}`));
	onClickOutside(popupBanRef, event => {
		if (!popupBan.value.visible) { return; }
		setTimeout(() => { popupBan.value.visible = false; }, 10);
	});

	// ----- Filters

	const isPostHidden = (): boolean => {
		const isFiltered = isPostFiltered();
		return isFiltered;
	}

	const isPostFiltered = (): boolean => {
		return (filterApplied.value?.action == FilterAction.Hide);
	}
	
	const refreshFilter = () => {
		const filters = LoadFilters();
		filterApplied.value = GetFilterMatchingPost(props.board, props.thread, props.post, filters);
	}

	const isPostFilterHighlighted = (): boolean => {
		return (filterApplied.value?.action == FilterAction.Highlight);
	}
</script>

<template>
	<div :id="`p${post.num}`" class="postContainer">
		<div v-if="!peek && popupDelete.visible" ref="delete-popup" class="popup popup-menu">
			<a v-if="post.filename" href="#" @click.prevent="deletePostFile()">Delete file</a>
			<a href="#" @click.prevent="deletePost()">Delete
				<template v-if="post.thread_num == post.num"><b>thread</b></template>
				<template v-else>post</template>
			</a>
		</div>

		<div v-if="!peek && popupBan.visible" ref="ban-popup" class="popup popup-ban">
			<CreateBanComponent :post="undefined" :currentBoard="board" />
		</div>

		<span v-if="!peek && thread.post_num != post.num" class="sideArrows"> &gt;&gt; </span>
		<span class="post" :class="{
			highlightedPost: is_highlighted,
			opPost: is_op_post,
			lastSeenPost: is_last_seen,
			notOP: !is_op_post,
			hidden: isPostHidden() && !peek,
		}" :style="isPostFilterHighlighted() ? `border: 3px solid ${filterApplied?.colorHex};` : ``">

		<div class="post-header" >
			<!-- Thread management by staff user -->
			<div v-if="!peek && userRole != undefined && post.num == post.thread_num && !thread.archived" class="admin-tools">
				<button v-if="canStickyThread()"  @click="onToggleSticky(thread)">Toggle Sticky <img src="/icons/sticky.png" /></button>
				<button v-if="canLockThread()"    @click="onToggleLocked(thread)">Toggle Locked <img src="/icons/lock.png" /></button>
				<button v-if="canArchiveThread()" @click="onArchive(thread)">Archive <img src="/icons/archive.png" /></button>
				<button v-if="canDeleteThread()"  @click="onDelete(thread)">Delete <img src="/icons/delete.png" /></button>
				<span v-if="canChangeThreadAutoCycle()">
					<label for=""><abbr title="Limit thread to this many posts.
Older posts are deleted.
A value of 0 (default) disables auto-cycle.">Auto-cycle</abbr>:</label>
					<input type="number" min="0" class="auto-cycle-input" v-model="thread.auto_cycle" />
					<button @click="onChangeAutoCycle(thread)">Set auto-cycle <img src="/icons/refresh.png" /></button>
				</span>
			</div>

			<span class="subject" v-if="thread.subject && thread.post_num == post.num">{{ thread.subject }}</span>
			<span class="username">{{ post.name ? post.name : "Anonymous" }}</span>
			<span class="capcode" :class="post.user_role" v-if="post.capcode">
				## {{ post.user_role }}
				<img :src="`/icons/user-role-${post.user_role}.gif`" class="icon" :title="post.user_role" />
			</span>
			<span class="tripcode" v-if="post.tripcode">{{ post.tripcode }}</span>
			<span
				v-if="board.config.ids_enabled && post.public_id"
				class="publicId"
				:style="{backgroundColor: GetPublicIdColorBackground(post.public_id), color: GetPublicIdColorForeground(post.public_id)}"
				@click="onClickPublicId(post.public_id)">
				ID:{{post.public_id}}
				<template v-if="public_id_count"> ({{public_id_count[post.public_id]}}) </template>
			</span>
			<!-- Post management by staff user -->
			<span v-if="userRole != undefined && !thread.archived" class="admin-tools-small">
				<a href="#" @click.prevent="clickBan(post)" title="Ban user" :id="`ban-post-btn-${id}`">
					<img src="/icons/no.png" alt="Ban" />
				</a>

				<a href="#" @click.prevent="clickDelete(post)" title="Delete post" :id="`delete-post-btn-${id}`">
					<img src="/icons/delete.png" alt="Delete"/>
				</a>

				<a href="#" @click.prevent="clickEdit(post)" title="Edit post">
					<img src="/icons/edit.png" alt="Edit"/>
				</a>
			</span>
			<span class="date">{{ GetPostTimeReadable(post.created_at) }}</span>
			<span class="postno"><a @click.prevent="onClickPostNo(post.num)" href="#" class="postNumLink">No.</a></span>
			<span class="postnum"><a @click.prevent="onClickPostNumber(post.num)" href="#" class="postNumLink">{{ post.num }}</a></span>
			<span v-if="post.num == post.thread_num">
				<img src="/icons/sticky.png" v-if="thread.sticky" title="Sticky"/>
				<img src="/icons/lock.png" v-if="thread.locked" title="Locked"/>
				<img src="/icons/archive.png" v-if="thread.archived" :title="`Archived at ${GetPostTimeReadable(thread.archived_at)}`"/>
				<img src="/icons/refresh.png" v-if="thread.auto_cycle > 0" :title="`Auto-Cycle`"/>
			</span>
			<span v-if="post.deleted_at != null">
				<img src="/icons/trash.png" :title="`Deleted at ${GetPostTimeReadable(post.deleted_at)}`"/>
			</span>
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

		<div class="post-body-container" v-if="!(isPostHidden() && !peek)">
			<span v-if="post.filename" class="post-file-container">
				<div class="post-file-container-header">
					File: <a :href="CdnAPI.GetPostImageURI(post)" target="_blank">{{ post.src_filename }}</a>
					({{GetFileSizeByteString(post.filesize)}}, {{post.img_width}}x{{post.img_height}})
					<template v-if="image_data?.expanded && isPostVideo(post)">
						- <a href="#" @click.prevent="onClickPostImage(post.id)">[Close]</a>
					</template>
				</div>

				<a :href="CdnAPI.GetPostImageURI(post)" target="_blank" @click.prevent class="post-file-link">
					<!-- Thumbnail -->

					<template v-if="!image_data?.expanded">
						<!-- Spoiler -->
						<img v-if="board.config.allow_spoilers && post.spoiler"
							@click="onClickPostImage(post.id)"
							class="post-image-thumb spoiler"
							:src="CdnAPI.GetSpoilerURI(board.config.spoiler_image)"
						>
						<!-- Regular thumbnail -->
						<img v-else
							@click="onClickPostImage(post.id)"
							class="post-image-thumb"
							:src="CdnAPI.GetPostImageThumbnailURI(post)"
						>
					</template>

					<!-- Real image/video -->

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
				<template v-if="post.md5 && post.thread_num == post.num">
					<img :src="CdnAPI.GetPublicURI('file_deleted.png')" class="post-image-thumb" title="File deleted" />
				</template>
			</div>

			<span class="post-body">
				<span v-if="isEditingPost">
					<textarea v-model="editingPostContent" cols="40" rows="6" placeholder="Post content"></textarea>
					<br/>
					<textarea v-model="editingPostHtml" cols="40" rows="6" placeholder="Post HTML"></textarea>
					<br/>
					<button @click="submitEditingPost">Submit changes</button>
				</span>
				<span v-else>
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
							<template v-else-if="token.type == 'spoiler'">
								<span class="spoiler-text">{{ token.text }}</span>
							</template>
						</span>
					</span>
					<span v-if="post.html">
						<br/>
						<span v-html="post.html"></span>
					</span>
				</span>
			</span>
		</div>
	</span>
</div>
</template>

<style scoped>
	.popup {
		z-index: 10;
		position: absolute;
		border: 1px solid var(--post-border-color);
		background-color: var(--post-background-color);

	}

	.popup-ban {
		padding: 1em;
		box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
	}

	.popup-menu {
		* {
			border-bottom: 1px solid var(--background-color-accent);
			padding: 0.2em;
		}

		*:hover {
			background-color: var(--background-color);
		}

		a {
			text-decoration: none;
			display: block;
		}
	}

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

			&.hidden {
				padding-bottom: 0.2em;
			}

			&.opPost {
				background-color: var(--background-color) !important;
			}

			&.notOP {
				border: 1px solid var(--post-border-color);
			}

			&.highlightedPost {
				&.notOP {
					border: 1px solid var(--highlighted-post-border-color);
				}
				background-color: var(--highlighted-post-background-color);
			}

			.post-header {
				span {
					margin-right: 0.25em;
				}

				.admin-tools {
					margin: 1em;
					padding: 0.5em;
					background-color: var(--background-color-darker);
					border: 1px solid var(--highlighted-post-border-color);

					>* {
						margin-right: 5px;
					}

					.auto-cycle-input {
						max-width: 4em;
					}
				}

				.admin-tools-small {
					>* {
						padding: 4px 4px 1px 4px;
						margin: 0 2px;
						border: 1px solid var(--highlighted-post-border-color);
					}

					>*:hover {
						background-color: var(--background-color-darker);
						border-color: var(--link-hover-color);
					}
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

				.publicId {
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

							&.spoiler {
								display: inline;
								width: 160px !important;
								height: 160px !important;
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
	}

	.error {
		color: var(--user-error-color);
	}

	.strikethrough {
		text-decoration: line-through;
	}
</style>
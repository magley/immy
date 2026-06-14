<script setup lang="ts">
	import type { BoardDTO } from '@/api/board.api';
	import { CdnAPI } from '@/api/cdn.api';
	import type { PostDTO } from '@/api/post.api';
	import type { ThreadDTO } from '@/api/thread.api';
import { GetMimeTypeFromFilename } from '@/util/file.util';
import { clamp } from '@vueuse/core';
	import { onMounted, onUnmounted, ref } from 'vue';

	interface GalleryModeProps {
		board: BoardDTO;
		thread: ThreadDTO;
		posts: PostDTO[];
	}

	const open = ref<boolean>(false);

	const props = defineProps<GalleryModeProps>();
	const currentPost = ref<PostDTO | undefined>(undefined);
	const currentPostIndex = ref<number>(0);

	const muted = ref<boolean>(false);
	const volume = ref<number>(0.5);

	const postsWithAnImage = (): PostDTO[] => {
		return props.posts.filter((p: PostDTO) => p.filename != "");
	}

	const prevImage = () => {
		currentPostIndex.value -= 1;
		if (currentPostIndex.value < 0) {
			currentPostIndex.value = 0;
		}
		currentPostIndexChanged();
	}

	const nextImage = () => {
		currentPostIndex.value += 1;
		const MAX: number = postsWithAnImage().length;
		if (currentPostIndex.value >= MAX) {
			currentPostIndex.value = MAX - 1;
		}
		currentPostIndexChanged();
	}

	const exitGalleryMode = () => {
		open.value = false;
	}

	const OpenGalleryMode = () => {
		open.value = true;
		currentPostIndexChanged();
	}

	const currentPostIndexChanged = () => {
		// (document.getElementById("gallery-video") as HTMLVideoElement).volume = volume.value;
		// (document.getElementById("gallery-video") as HTMLVideoElement).muted = muted.value;
		currentPost.value = postsWithAnImage()[currentPostIndex.value];
	}

	const imageCount = (): number => {
		return postsWithAnImage().length;
	}

	onMounted(() => {
		window.addEventListener("keydown", onKey);
	});

	onUnmounted(() => {
		window.removeEventListener("keydown", onKey);
	});

	const onKey = (e: KeyboardEvent) => {
		if (e.key == "ArrowLeft") {
			prevImage();
		} else if (e.key == "ArrowRight") {
			nextImage();
		} else if (e.key == "Escape") {
			exitGalleryMode();
		} else if (e.key == "m") {
			muted.value = !muted.value;
		} else if (e.key == ",") {
			volume.value -= 0.05;
			volume.value = clamp(volume.value, 0, 1);
		} else if (e.key == ".") {
			volume.value += 0.05;
			volume.value = clamp(volume.value, 0, 1);
		}
	}

	defineExpose({OpenGalleryMode});

	const isPostImage = (post: PostDTO) => {
		return GetMimeTypeFromFilename(post.filename).startsWith('image/');
	}
	const isPostVideo = (post: PostDTO) => {
		return GetMimeTypeFromFilename(post.filename).startsWith('video/');
	}

</script>

<template>
	<div id="gallery-pane" v-if="open">
		<div class="gallery-header">
			<div title="M">
				<label for="toolbar-muted"><img src="/icons/sound-off.png" v-show="muted" /><img src="/icons/sound.png" v-show="!muted" />Muted:</label>
				<input type="checkbox" v-model="muted" id="toolbar-muted" />
			</div>

			<div title=", - lower volume
. - raise volume">
				<label for="toolbar-volume"><img src="/icons/volume.png" />Volume:</label>
				<input type="range" v-model="volume" id="toolbar-volume" :min="0" :max="1" :step="0.05" />
			</div>

			<button @click="prevImage" title="Left Arrow Key" :disabled="currentPostIndex == 0">
				<img src="/icons/back.png" /> Back
			</button>
			<button @click="nextImage" title="Right Arrow Key" :disabled="currentPostIndex == imageCount() - 1">
				<img src="/icons/forward.png" /> Next
			</button>
			<button @click="exitGalleryMode" title="Escape" >
				<img src="/icons/cancel.png" /> Exit
			</button>
		</div>
		<div class="gallery-image-container" v-if="currentPost">
			<img v-if="isPostImage(currentPost)"
				:src="CdnAPI.GetPostImageURI(currentPost)"
				class="gallery-image" />
			<video v-if="isPostVideo(currentPost)"
				id="gallery-video"
				:key="currentPost.filename"
				class="gallery-image"
				controls
				autoplay
				:muted="muted"
				:volume="volume"
				:width="currentPost.img_width"
				:height="currentPost.img_height"
				loop>
				<source :src="CdnAPI.GetPostImageURI(currentPost)"
				:type="GetMimeTypeFromFilename(currentPost.filename)" />
				Your browser does not support the video tag.
			</video>
		</div>
	</div>
</template>

<style scoped>
	#gallery-pane {
		position: absolute;
		top: 0;
		left: 0;
		width: 100vw;
		height: 100vh;
		z-index: 1000;
		background-color: rgba(0, 0, 0, 0.9);

		.gallery-header {
			width: 100%;
			display: flex;
			gap: 1em;
			font-size: 1.75em;
			font-family: monospace;

			justify-content: center;
			color: white;
		}

		.gallery-image-container {
			display: flex;
			justify-content: center;

			width: 90%;
			height: 90%;
			margin: auto;

			.gallery-image {
				max-width: 100%;
				max-height: 100%;
			}
		}
	}

	input[type=range] {
    	margin: 0;
 		height: 1em;
 	}
</style>
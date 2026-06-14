<script setup lang="ts">
	import type { BoardDTO } from '@/api/board.api';
	import { CdnAPI } from '@/api/cdn.api';
	import type { PostDTO } from '@/api/post.api';
	import type { ThreadDTO } from '@/api/thread.api';
	import { onMounted, onUnmounted, ref } from 'vue';

	interface GalleryModeProps {
		board: BoardDTO;
		thread: ThreadDTO;
		posts: PostDTO[];
	}

	const open = ref<boolean>(false);

	const props = defineProps<GalleryModeProps>();
	const currentPostIndex = ref<number>(0);

	const postsWithAnImage = (): PostDTO[] => {
		return props.posts.filter((p: PostDTO) => p.filename != "");
	}

	const currentPost = (): PostDTO | undefined => {
		return postsWithAnImage()[currentPostIndex.value];
	}

	const prevImage = () => {
		currentPostIndex.value -= 1;
		if (currentPostIndex.value < 0) {
			currentPostIndex.value = 0;
		}
	}

	const nextImage = () => {
		currentPostIndex.value += 1;
		const MAX: number = postsWithAnImage().length;
		if (currentPostIndex.value >= MAX) {
			currentPostIndex.value = MAX - 1;
		}
	}

	const exitGalleryMode = () => {
		open.value = false;
	}

	const OpenGalleryMode = () => {
		open.value = true;
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
		}
	}

	defineExpose({OpenGalleryMode});
</script>

<template>
	<div id="gallery-pane" v-if="open">
		<div class="gallery-header">
			<button @click="prevImage" title="Left Arrow Key" :disabled="currentPostIndex == 0">Back</button>
			<span>{{ currentPostIndex + 1}} / {{ imageCount() }} </span>
			<button @click="nextImage" title="Right Arrow Key" :disabled="currentPostIndex == imageCount() - 1">Next</button>
			<button @click="exitGalleryMode" title="Escape" >Exit</button>
		</div>
		<div class="gallery-image-container">
			<img
				class="gallery-image"
				v-if="currentPost()"
				:src="CdnAPI.GetPostImageURI(currentPost()!)"
				/>
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
			display: flex;
			gap: 1em;

			justify-content: center;
			font-size: 2em;
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
</style>
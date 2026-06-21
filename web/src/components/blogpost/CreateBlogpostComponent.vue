<script setup lang="ts">
	import { type CreateBlogpostDTO, type BlogpostDTO, BlogpostAPI } from '@/api/blogpost.api';
	import type { AxiosError } from 'axios';
	import { ref } from 'vue';

	const error = ref<string | undefined>(undefined);
	const preview = ref<boolean>(false);
	const emits = defineEmits(["createdBlogpost"]);

	const dto = ref<CreateBlogpostDTO>({
		title: 'Untitled Blogpost',
		html: 'This is my <b>blogpost</b>!'
	});

	const createBlogpost = () => {
		error.value = undefined;

		BlogpostAPI.CreateBlogpost(dto.value).then((res) => {
			dto.value.title = "";
			dto.value.html = "";
			emits("createdBlogpost");
		}).catch((err: AxiosError) => {
			if (err.response?.status == 401) {
				error.value = "You are not authorized to created blogposts."
			} else {
				error.value = "Could not create blogpost";
			}
			console.error(err);
		});
	}
</script>

<template>
	<div class="container">
		<div class="center form">
			<label for="blog-title">Title:</label><br/>
			<input id="blog-title" placeholder="Enter your title" v-model="dto.title" />
			<br/>
			<label for="blog-html">Message:</label><br/>
			<textarea id="blog-html" placeholder="Write your blogpost here (HTML supported)" cols="50" rows="10" v-model="dto.html" />
			<br/>
			<button @click="createBlogpost">Create</button>
			<br/>
			<div class="error" v-if="error">{{error}}</div>
		</div>
		<div>
			<i>(Preview)</i>
			<div class="center preview">
				<h2>{{dto.title}}</h2>
				<p v-html="dto.html"></p>
			</div>
		</div>
	</div>
</template>

<style scoped>
	.container {
		display: flex;
		align-items: flex-start;

		>div {
			flex: 1;
			margin: 0 1em;
		}

		width: 86%;
		margin: auto;
	}

	.error {
		color: var(--user-error-color);
	}

	.form {
		input, textarea {
			width: 100%;
		}
	}

	.center {
		margin: auto;
	}

	.preview {
		border: 1px solid var(--banner-title-color);
		background: white;
		padding: 1em;

		h2 {
			border: 1px solid var(--banner-title-color);
			background: var(--background-color);
			color: var(--banner-title-color);
			padding: 0.25em;
		}

		p {
			white-space: pre-wrap;
		}
	}
</style>
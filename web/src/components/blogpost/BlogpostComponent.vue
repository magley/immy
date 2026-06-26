<script setup lang="ts">
	import type { BlogpostDTO } from '@/api/blogpost.api';
import { GetPostTimeReadable } from '@/model/post/post.model';

	interface BlogpostProps {
		blogpost: BlogpostDTO
	}

	const props = defineProps<BlogpostProps>();
</script>

<template>
	<div class="blogpost">
		<h2>
			<RouterLink :to="`/blog/${blogpost.id}`">{{ blogpost.title }}</RouterLink>
			by {{ blogpost.author_name }}
			@ {{ GetPostTimeReadable(blogpost.created_at) }}
		</h2>
		<p v-html="blogpost.html"></p>
	</div>
</template>

<style scoped>
	.blogpost {
		width: 60%;
		margin-left: auto;
		margin-right: auto;

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
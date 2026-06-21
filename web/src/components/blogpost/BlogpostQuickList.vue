<script setup lang="ts">
	import { BlogpostAPI, type BlogpostShortDTO } from '@/api/blogpost.api';
	import type { AxiosError } from 'axios';
	import { onMounted, ref } from 'vue';

	const blogposts = ref<BlogpostShortDTO[]>([]);
	const error = ref<string | undefined>(undefined);

	onMounted(() => {
		BlogpostAPI.ListBlogpostsShort(0, 5).then((res) => {
			blogposts.value = res.data.data!;
		}).catch((err: AxiosError) => {
			error.value = "Could not fetch blogposts";
			console.error(err);
		});
	})

	const toDate = (blog: BlogpostShortDTO): string => {
		const date = new Date(Date.parse(blog.created_at));
		return `${date.getDate()}/${date.getMonth()}/${date.getFullYear() - 2000}`;
	}
</script>

<template>
	<div class="blogpost-quick-list">
		<hr/>
		<div v-if="error" class="error">
			{{error}}
		</div>
		<div v-else>
			<div v-for="blog of blogposts">
				{{ toDate(blog) }} <RouterLink :to="`blog/${blog.id}`">{{ blog.title }}</RouterLink>
			</div>
		</div>
		<hr/>
	</div>
</template>

<style scoped>
	.blogpost-quick-list {
		width: 468px; /* Board banner image width */
		margin: auto;
	}

	.error {
		color: var(--user-error-color);
	}
</style>
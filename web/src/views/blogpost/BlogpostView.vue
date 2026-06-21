<script setup lang="ts">
	import { BlogpostAPI, type BlogpostDTO } from '@/api/blogpost.api';
	import BlogpostComponent from '@/components/blogpost/BlogpostComponent.vue';
	import CreateBlogpostComponent from '@/components/blogpost/CreateBlogpostComponent.vue';
	import type { AxiosError } from 'axios';
	import { onMounted, ref } from 'vue';

	const blogposts = ref<BlogpostDTO[]>([]);

	const loading = ref<boolean>(false);
	const error = ref<string | undefined>(undefined);


	const perPage = 20;
	const page = ref<number>(1);
	const pagesTotal = ref<number>(1);
	const pagesNav = ref<number[]>([]);

	onMounted(() => {
		getBlogposts();
	})

	const getBlogposts = () => {
		if (page.value < 1) page.value = 1;
		if (page.value > pagesTotal.value) page.value = pagesTotal.value;

		loading.value = true;

		BlogpostAPI.ListBlogposts((page.value - 1) * perPage, perPage).then((res) => {
			blogposts.value = res.data.data!;

			const meta = res.data.meta!;
			page.value = meta.page;
			pagesTotal.value = meta.total_pages;
			pagesNav.value = [
				page.value - 4, page.value - 3, page.value - 2, page.value - 1,
				page.value - 0,
				page.value + 1, page.value + 2, page.value + 3, page.value + 4, page.value + 5,
			];
			pagesNav.value = pagesNav.value.filter((v) => v >= 1 && v <= meta.total_pages);
		}).catch((err: AxiosError) => {
			error.value = "Could not fetch blogposts";
			console.error(err);
		}).finally(() => {loading.value = false;});
	}

	const gotoPage = (p: number) => {
		page.value = p;
		getBlogposts();
	}

	const onBlogpostCreated = () => {
		getBlogposts();
	}
</script>

<template>
	<h1>Blogposts</h1>
	<div class="error" v-if="error">{{ error }}</div>


	<template v-if="blogposts.length == 0">
		<div class="center">No blogposts have been written yet.</div>
	</template>
	<template v-else>
		<BlogpostComponent v-for="blogpost of blogposts" :blogpost="blogpost" />

		<div class="center nav">
			[<a href="#" @click.prevent="gotoPage(1)">First</a>]&thinsp;
			[<a href="#" @click.prevent="gotoPage(page - 1)">Prev</a>]&thinsp;
			<span v-for="p, i of pagesNav">
				<template v-if="p == page">
					<span>{{ p }} </span>
				</template>
				<template v-else>
					<a href="#" @click.prevent="gotoPage(p)">{{ p }} </a>
				</template>
				<template v-if="i < pagesNav.length - 1">,</template>&thinsp;
			</span>
			[<a href="#" @click.prevent="gotoPage(page + 1)">Next</a>]&thinsp;
			[<a href="#" @click.prevent="gotoPage(pagesTotal)">Last</a>]&thinsp;
		</div>
	</template>


</template>

<style scoped>
	.error {
		color: var(--user-error-color);
	}

	.center {
		text-align: center;
	}

	h1, h2 {
		text-align: center;
		color: var(--banner-title-color);
	}
</style>
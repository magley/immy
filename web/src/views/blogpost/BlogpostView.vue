<script setup lang="ts">
	import { BlogpostAPI, type BlogpostDTO } from '@/api/blogpost.api';
	import BlogpostComponent from '@/components/blogpost/BlogpostComponent.vue';
	import CreateBlogpostComponent from '@/components/blogpost/CreateBlogpostComponent.vue';
	import PaginatorComponent from '@/components/PaginatorComponent.vue';
	import { Paginator } from '@/util/pagination.util';
	import { onMounted, reactive, ref } from 'vue';

	const paginator = reactive<Paginator<BlogpostDTO[]>>(new Paginator(BlogpostAPI.ListBlogposts));
	const blogposts = ref<BlogpostDTO[]>([]);
	const error = ref<string | undefined>(undefined);
	const collapsedNewBlogpostForm = ref<boolean>(true);
	const userRole = ref<string | undefined>(undefined);

	onMounted(() => {
		getBlogposts();
		userRole.value = localStorage.getItem("role") ?? undefined;
	})

	const getBlogposts = () => {
		paginator.getItems()
			.then((res) => blogposts.value = res.data.data!)
			.catch((_) => error.value = "Could not fetch blogposts!");
	}

	const gotoPage = (p: number) => {
		paginator.page = p;
		getBlogposts();
	}

	const onBlogpostCreated = () => {
		getBlogposts();
	}
</script>

<template>
	<h1>Blogposts</h1>
	<div class="error" v-if="error">{{ error }}</div>

	<div v-if="userRole == 'admin'">
		<hr />
		<div v-if="collapsedNewBlogpostForm" class="center">
			<br/>
			[ <a href="#" @click.prevent="() => {collapsedNewBlogpostForm = false;}">Create new blogpost</a> ]
			<br/><br/>
		</div>
		<div v-else>
			<div class="center">
				<br/>
				[ <a href="#" @click.prevent="() => {collapsedNewBlogpostForm = true;}">Hide</a> ]
			</div>
			<h2>New blogpost</h2>
			<CreateBlogpostComponent :blogpost-to-edit="undefined" @created-blogpost="onBlogpostCreated" />
		</div>
		<hr />
	</div>

	<template v-if="blogposts.length == 0">
		<div class="center">No blogposts have been written yet.</div>
	</template>
	<template v-else>
		<PaginatorComponent :paginator="paginator" @goto-page="gotoPage" />
		<BlogpostComponent v-for="blogpost of blogposts" :blogpost="blogpost" />
		<PaginatorComponent :paginator="paginator" @goto-page="gotoPage" />
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
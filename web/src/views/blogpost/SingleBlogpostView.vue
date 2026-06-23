<script setup lang="ts">
	import { BlogpostAPI, type BlogpostDTO } from '@/api/blogpost.api';
	import BlogpostComponent from '@/components/blogpost/BlogpostComponent.vue';
	import CreateBlogpostComponent from '@/components/blogpost/CreateBlogpostComponent.vue';
	import type { AxiosError } from 'axios';
	import { onMounted, ref } from 'vue';
	import { useRoute, useRouter } from 'vue-router';

	const route = useRoute();
	const router = useRouter();

	const blogpost = ref<BlogpostDTO | undefined>(undefined);

	const loading = ref<boolean>(false);
	const error = ref<string | undefined>(undefined);

	const userRole = ref<string | undefined>(undefined);

	const editMode = ref<boolean>(false);

	onMounted(() => {

		getBlogpost();
		userRole.value = localStorage.getItem("role") ?? undefined;
	})

	const getBlogpost = () => {
		loading.value = true;

		const idStr = route.params.id;
		if (idStr == undefined) {
			router.push("/blog");
			return;
		}

		const id = Number(idStr);
		if (isNaN(id)) {
			router.push("/blog");
			return;
		}

		BlogpostAPI.GetBlogpost(id).then((res) => {
			blogpost.value = res.data.data!;
		}).catch((err: AxiosError) => {
			error.value = "Could not fetch blogpost";
			console.error(err);
		}).finally(() => {loading.value = false;});
	}

	const toggleEditMode = () => {
		editMode.value = !editMode.value;
	}

	const onUpdatedBlogpost = (newBlogpost: BlogpostDTO) => {
		blogpost.value = newBlogpost;
		editMode.value = false;
	}
</script>

<template>
	<h1>Blogpost <template v-if="blogpost">#{{blogpost.id}}</template></h1>
	<div class="error" v-if="error">
		{{ error }}
	</div>

	<template v-if="blogpost">
		<h3 v-if="userRole == 'admin'" class="center narrow">
			[ <a href="#" @click.prevent="toggleEditMode">Toggle edit mode</a> ]
		</h3>
		<BlogpostComponent :blogpost="blogpost" v-if="!editMode" />
		<CreateBlogpostComponent :blogpost-to-edit="blogpost" v-if="editMode" @updated-blogpost="onUpdatedBlogpost" />
	</template>

	<div class = "narrow center">
		<hr />
		Return to <RouterLink to="/blog">blogposts</RouterLink>.
	</div>
</template>

<style scoped>
	.error {
		color: var(--user-error-color);
	}

	.center {
		text-align: center;
	}

	.narrow {
		width: 54%;
		margin: auto;
		margin-top: 1em;
	}

	h1, h2 {
		text-align: center;
		color: var(--banner-title-color);
	}
</style>
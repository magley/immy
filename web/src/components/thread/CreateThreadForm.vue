<script setup lang="ts">
	import type { ApiResponse } from '@/api/http';
	import { type CreatePostForThreadDTO } from "@/api/post.api.ts";
	import { ThreadAPI, type CreateThreadDTO, type ThreadDTO } from '@/api/thread.api';
	import { FileToBase64, GetFileFromEvent } from '@/util/file.util';
	import type { AxiosError, AxiosResponse } from 'axios';
	import { ref } from 'vue';

	export interface CreateThreadProps {
		board_code: string,
		max_size_bytes: number,
	}

	const props = defineProps<CreateThreadProps>();
	const emit = defineEmits(['threadCreated']);

	const createThreadDTO = ref<CreateThreadDTO>({
		board_code: '',
		subject: '',
		locked: undefined,
		sticky: undefined,
		post: undefined
	});
	const createThreadPostDTO = ref<CreatePostForThreadDTO>({
		name: '',
		content: '',
		filename: '',
		filebytes: '',
		options: ''
	});


	const createThreadError = ref<string | undefined>(undefined);

	const fileName = ref<string | null>(null);
	const fileBytes = ref<string | null>(null);
	const fileError = ref<string | null>(null);

	const onSubmitReply = () => {
		createThreadError.value = undefined;

		if (!fileName.value || !fileBytes.value) {
			createThreadError.value = "No file attached";
			return;
		}

		createThreadPostDTO.value.filename = fileName.value;
		createThreadPostDTO.value.filebytes = fileBytes.value;

		createThreadDTO.value.board_code = props.board_code;
		createThreadDTO.value.post = createThreadPostDTO.value;

		ThreadAPI.CreateThread(createThreadDTO.value).then((res: AxiosResponse<ApiResponse<ThreadDTO>>) => {
			emit('threadCreated');
		}).catch((err: AxiosError) => {
			createThreadError.value = "Could not create thread";
			console.error(err);
		});
	}

	const onFileSelected = async (e: Event) => {
		fileError.value = null;

		const [file, err] = GetFileFromEvent(e, props.max_size_bytes);

		if (err) {
			fileError.value = err;
		}
		if (!file) {
			clearSelectedFile();
			return;
		}

		FileToBase64(file).then((b64: string) => {
			fileName.value = file.name;
			fileBytes.value = b64;
		}).catch((err: any) => {
			fileError.value = "Could not process file";
			console.error(err);
			clearSelectedFile();
			return;
		});
	}

	const clearSelectedFile = () => {
		const fileInput: HTMLFormElement = document.getElementById("reply-file-upload") as HTMLFormElement;
		fileInput.value = null;
		fileBytes.value = null;
		fileName.value = null;
	}
</script>

<template>
	<form @submit.prevent="onSubmitReply">
		<input type=text placeholder="Subject" v-model="createThreadDTO.subject"/><br/>
		<input type=text placeholder="Name" v-model="createThreadPostDTO.name"/><br/>
		<input type=text placeholder="Options" v-model="createThreadPostDTO.options"/><br/>
		<textarea cols=30 rows=10 id="reply-area" placeholder="Text..." v-model="createThreadPostDTO.content"/><br/>
		<input type="file" accept="image/png, image/jpeg, image/gif" @change="onFileSelected" id="reply-file-upload"><br/>
		<template v-if="fileError"><span class="error">{{fileError}}</span></template>

		<br/>
		<button type=submit>Post thread</button>

		<template v-if="createThreadError">
			<div/>
			<span class="error">{{createThreadError}}</span>
		</template>
	</form>
</template>

<style scoped>
	.error {
		color: red;
	}
</style>
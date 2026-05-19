<script setup lang="ts">
	import { ref, useTemplateRef } from 'vue';
	import { PostAPI, type CreatePostDTO, type PostDTO } from "@/api/post.api.ts";
	import type { ApiResponse } from '@/api/http';
	import type { AxiosResponse, AxiosError } from 'axios';
	import { FileToBase64, GetFileFromEvent } from '@/util/file.util';
	
	export interface CreatePostProps {
		thread_id: number,
		max_size_bytes: number,
	}

	const props = defineProps<CreatePostProps>();
	const emit = defineEmits(['postCreated']);

	const replyDTO = ref<CreatePostDTO>({
		name: '',
		content: '',
		filename: null,
		filebytes: null,
		options: '',
		thread_id: 0
	});
	const replyError = ref<string | undefined>(undefined);
	
	const fileName = ref<string | null>(null);
	const fileBytes = ref<string | null>(null);
	const fileError = ref<string | null>(null);

	const textArea = useTemplateRef('text-area');

	const onSubmitReply = () => {
		replyError.value = undefined;

		replyDTO.value.thread_id = props.thread_id;
		replyDTO.value.filename = fileName.value ?? null;
		replyDTO.value.filebytes = fileBytes.value ?? null;

		PostAPI.CreatePost(replyDTO.value).then((res: AxiosResponse<ApiResponse<PostDTO>>) => {
			emit('postCreated');
		}).catch((err: AxiosError) => {
			replyError.value = "Could not post reply";
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
			console.log(b64);
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

	const AppendText = (text: string) => {
		if (replyDTO.value.content == undefined) {
			replyDTO.value.content = "";
		}
		replyDTO.value.content += text;

		if (!textArea.value) {
			console.error("Text area is null. This shouldn't happen");
			return;
		}
		textArea.value.scrollTop = textArea.value.scrollHeight;
	}

	defineExpose({ AppendText });
</script>

<template>
	<form @submit.prevent="onSubmitReply">
		<input type=text placeholder="Name" v-model="replyDTO.name"/><br/>
		<input type=text placeholder="Options" v-model="replyDTO.options"/><br/>
		<textarea id="reply-area" placeholder="Text..." ref='text-area' v-model="replyDTO.content"/><br/>
		<input type="file" accept="image/png, image/jpeg" @change="onFileSelected" id="reply-file-upload"><br/>
		<template v-if="fileError"><span class="error">{{fileError}}</span></template>

		<br/>
		<button type=submit>Post reply</button>

		<template v-if="replyError">
			<div/>
			<span class="error">{{replyError}}</span>
		</template>
	</form>
</template>

<style scoped>
	.error {
		color: red;
	}
</style>
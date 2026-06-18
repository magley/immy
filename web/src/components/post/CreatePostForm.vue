<script setup lang="ts">
	import { ref, useTemplateRef } from 'vue';
	import { PostAPI, type CreatePostDTO, type PostDTO, type CreatePostForThreadDTO } from "@/api/post.api.ts";
	import type { ApiErrorInfo, ApiResponse } from '@/api/http';
	import type { AxiosResponse, AxiosError } from 'axios';
	import { FileToBase64, GetFileFromEvent } from '@/util/file.util';
	import { type BoardDTO } from '@/api/board.api';
	import { ThreadAPI, type CreateThreadDTO, type ThreadDTO } from '@/api/thread.api';
	import { useId } from 'vue';

	const id = useId();

	export interface CreatePostProps {
		/** If `thread` is undefined, then it's assumed that this
		  * post is being created as part of a new thread. */
		thread: ThreadDTO | undefined,
		max_size_bytes: number,
		mime_types_allowed: string[],
		board: BoardDTO,
	}

	const props = defineProps<CreatePostProps>();
	const emit = defineEmits(['postCreated']);

	const replyDTO = ref<CreatePostDTO>({
		name: '',
		content: '',
		filename: null,
		filebytes: null,
		options: '',
		thread_id: 0,
		spoiler: false,
	});
	const replyError = ref<string | undefined>(undefined);
	
	const fileName = ref<string | null>(null);
	const fileBytes = ref<string | null>(null);
	const fileError = ref<string | null>(null);

	const subject = ref<string>('');
	const textArea = useTemplateRef('text-area');

	const spoilerExpanded = ref<boolean>(false);
	const mathExpanded = ref<boolean>(false);
	const codeExpanded = ref<boolean>(false);

	const Clear = () => {
		clearSelectedFile();
		replyDTO.value = {
			name: '',
			content: '',
			filename: null,
			filebytes: null,
			options: '',
			thread_id: 0,
			spoiler: false,
		};
		clearSelectedFile();
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

	defineExpose({ AppendText, Clear });

	const onSubmitReply = () => {
		if (isFormForNewThread()) {
			createThread();
		} else {
			createReply();
		}
	}

	const isFormForNewThread = () => props.thread == undefined;

	const createReply = () => {
		if (props.thread == undefined) return;
		if (props.thread.locked) return;

		replyError.value = undefined;

		replyDTO.value.thread_id = props.thread.id;
		replyDTO.value.filename = fileName.value ?? null;
		replyDTO.value.filebytes = fileBytes.value ?? null;

		PostAPI.CreatePost(replyDTO.value).then((res: AxiosResponse<ApiResponse<PostDTO>>) => {
			emit('postCreated');
		}).catch((err: AxiosError<ApiResponse<PostDTO>>) => {
			if (err.response?.data.error?.code == "BANNED") {
				replyError.value = "You are <a href='/banned'>banned</a>.";
				console.log("You are banned");
			} else if (err.response?.data.error?.code == "WARNED") {
				replyError.value = "You have been <a href='/banned'>warned</a>.";
			} else {
				replyError.value = "Could not post reply";
				console.error(err);
			}
		});
	}

	const createThread = () => {
		replyError.value = undefined;

		if (!fileName.value || !fileBytes.value) {
			replyError.value = "No file attached";
			return;
		}

		replyDTO.value.filename = fileName.value ?? null;
		replyDTO.value.filebytes = fileBytes.value ?? null;

		const createThreadPostDTO: CreatePostForThreadDTO = {
			name: replyDTO.value.name,
			content: replyDTO.value.content,
			filename: replyDTO.value.filename,
			filebytes: replyDTO.value.filebytes,
			options: replyDTO.value.options,
			spoiler: replyDTO.value.spoiler,
		};
		const createThreadDTO: CreateThreadDTO = {
			board_code: props.board.code,
			subject: subject.value,
			locked: false,
			sticky: false,
			post: createThreadPostDTO
		};

		ThreadAPI.CreateThread(createThreadDTO).then((res: AxiosResponse<ApiResponse<ThreadDTO>>) => {
			emit('postCreated');
		}).catch((err: AxiosError<ApiResponse<PostDTO>>) => {
			if (err.response?.data.error?.code == "BANNED") {
				replyError.value = "You are <a href='/banned'>banned</a>.";
			} else if (err.response?.data.error?.code == "WARNED") {
				replyError.value = "You have been <a href='/banned'>warned</a>.";
			} else {
				replyError.value = "Could not create thread";
				console.error(err);
			}
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
		const fileInput: HTMLFormElement = document.getElementById(`reply-file-upload-${id}`) as HTMLFormElement;
		fileInput.value = null;
		fileBytes.value = null;
		fileName.value = null;
	}

	const toggleHelpSpoiler = () => {
		spoilerExpanded.value = !spoilerExpanded.value;
	}
	const toggleHelpMath = () => {
		mathExpanded.value = !mathExpanded.value;
	}
	const toggleHelpCode = () => {
		codeExpanded.value = !codeExpanded.value;
	}
</script>

<template>
	<form @submit.prevent="onSubmitReply">
		<input type=text placeholder="Name" v-model="replyDTO.name"/><br/>
		<input type=text placeholder="Options" v-model="replyDTO.options"/><br/>
		<textarea cols=30 rows=10 :id="`reply-area-${id}`" placeholder="Text..." ref='text-area' v-model="replyDTO.content"/><br/>
		<input type="file" :accept="mime_types_allowed.join(', ')" @change="onFileSelected" :id="`reply-file-upload-${id}`"><br/>
		<template v-if="fileError"><span class="error">{{fileError}}</span></template>

		<div v-if="board.config.allow_spoilers" >
			<label :for="`create-post-spoiler-${id}`">
				Spoiler:
			</label>
			<input :id="`create-post-spoiler-${id}`" type=checkbox v-model="replyDTO.spoiler"/>
		</div>

		<br/>

		<div class="help-container">
			<div v-if="board.config.spoiler_image">
				Spoilers are enabled <a @click.prevent="toggleHelpSpoiler">[Help]</a>
				<div v-if="spoilerExpanded" class="help-explanation">
					Wrap hidden text with <tt>[spoiler] [/spoiler]</tt> tags:
<pre>
Darth Vader is [spoiler]Luke's father[/spoiler]
</pre>
					This will render:
					<p>
					Darth Vader is <span class="spoiler-text">Luke's father</span>
					</p>
				</div>
			</div>
			<div v-if="board.config.math_enabled">
				<vue-latex expression="\LaTeX" /> is supported <a @click.prevent="toggleHelpMath">[Help]</a>
				<div v-if="mathExpanded" class="help-explanation">
					Wrap your equations between <tt>[math] [/math]</tt> tags:
<pre>
The Pythagorean theorem is:
[math]
a^2 + b^2 = c^2
[/math]
</pre>
					This will render:
					<p>
					The Pythagorean theorem is:<br/>
					<vue-latex expression="a^2 + b^2 = c^2" />
					</p>
				</div>
			</div>
			<div v-if="board.config.code_enabled">
				<tt><strong>Syntax highlighting</strong></tt> is supported <a @click.prevent="toggleHelpCode">[Help]</a>
				<div v-if="codeExpanded" class="help-explanation">
					Wrap your code between <tt>[code] [/code]</tt> tags:
<pre>
Hello world in C:
[code]
#import &lt;stdio.h&gt;
int main() {
    printf("Hello world\n");
    return 0;
}
[/code]
</pre>
					This will render:
					<p>
					Hello world in C:<br/>
<highlightjs autodetect code='#import <stdio.h>
int main() {
    printf("Hello world\n");
    return 0;
}' />
					</p>
				</div>
			</div>
		</div>

		<button type=submit class="submit-button">
			<template v-if="isFormForNewThread()">Create new Thread</template>
			<template v-else>Post Reply</template>
		</button>

		<template v-if="replyError">
			<br/>
			<div class="postError" v-html="replyError"></div>
		</template>
	</form>
</template>

<style scoped>
	.error {
		color: var(--user-error-color);
	}

	.postError {
		color: var(--user-error-color);
		padding: 2px;
	}

	.help-container {
		a {
			cursor: pointer;
			user-select: none;
		}

		.help-explanation {
			display: block;
			background-color: var(--background-color-darker);
			border: 1px solid black;
			padding: 0.5em;

			text-align: left !important;
		}
	}

	.submit-button {
		margin: 1em 0;
		text-align: center;
		padding: 0.25em;
	}
</style>
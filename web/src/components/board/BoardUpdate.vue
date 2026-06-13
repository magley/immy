<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { BoardAPI, type BoardDTO, type UpdateBoardDTO, type CreateBoardDTO } from "@/api/board.api.ts";
	import type { AxiosError, AxiosResponse } from 'axios';
	import type { ApiResponse } from '@/api/http';
	import { GetElipsisString, IsAlphaNumeric, StripSlashes } from '@/util/various.util';
	import { GetFileSizeByteFromString, GetFileSizeByteString } from '@/util/file.util';
	import { MetaAPI } from '@/api/meta.api';
	import { CdnAPI } from '@/api/cdn.api';

	interface Props {
		currentBoardValue: BoardDTO | undefined;
		allowedMimeTypes: string[];
	}

	const props = defineProps<Props>();
	const emits = defineEmits(['created']);

	const spoiler_images = ref<string[]>([]);

	const boardDTO = ref<BoardDTO>({
		name: '',
		code: '',
		description: '',
		config: {
			locked: false,
			hidden: false,
			max_file_size: 0,
			reply_files_allowed: true,
			mime_types_allowed: [],
			bump_limit: 250,
			image_limit: 150,
			flags_enabled: false,
			ids_enabled: false,
			code_enabled: false,
			math_enabled: false,
			max_threads: 100,
			allow_spoilers: false,
			spoiler_image: "",
		},
		id: 0,
		created_at: '',
		deleted_at: '',
		meta: {
			post_count: 0,
			bytes_uploaded: 0
		}
	});

	const addMimeTypeInputVal = ref<string>("");
	const maxFileSizeInputVal = ref<string>("2 MB");
	const error = ref<string | undefined>(undefined);

	onMounted(() => {
		if (props.currentBoardValue) {
			boardDTO.value = Object.assign(boardDTO.value, props.currentBoardValue);
		}
		setMimePreset("image+video");

		CdnAPI.GetFilesIn("public/spoiler/").then((res) => {
			spoiler_images.value = res.sort();
			if (boardDTO.value && spoiler_images.value?.length > 0) {
				if (boardDTO.value.config.spoiler_image == "") {
					boardDTO.value.config.spoiler_image = spoiler_images.value[0]!;
				}
			}
		});
	});

	const isFormForNewBoard = () => props.currentBoardValue == undefined;

	const onSubmit = () => {
		error.value = "";
		if (!prevalidate()) {
			return;
		}

		if (isFormForNewBoard()) {
			createNewBoard();
		} else {
			updateThisBoard();
		}
	}

	const prevalidate = () => {
		boardDTO.value.config.max_file_size = GetFileSizeByteFromString(maxFileSizeInputVal.value);
		boardDTO.value.name = boardDTO.value.name.trim();
		boardDTO.value.code = StripSlashes(boardDTO.value.code.trim());
		boardDTO.value.description = boardDTO.value.description.trim();

		if (!IsAlphaNumeric(boardDTO.value.code)) {
			error.value = "Board code must consist of only alpha-numeric characters";
			return false;
		}
		if (isNaN(boardDTO.value.config.max_file_size)) {
			error.value = "Invalid maximum file size";
			return false;
		}
		if (boardDTO.value.config.mime_types_allowed.length == 0) {
			error.value = "No MIME types specified (users won't be able to attach files)";
			return false;
		}

		return true;
	}

	const createNewBoard = () => {
		const createBoardDTO: CreateBoardDTO = {
			name: boardDTO.value.name,
			code: boardDTO.value.code,
			description: boardDTO.value.description,
			config: boardDTO.value.config
		};

		BoardAPI.CreateBoard(createBoardDTO).then((res: AxiosResponse<ApiResponse<BoardDTO>>) => {
			emits("created", res.data.data!.id);
		}).catch((err: AxiosError) => {
			error.value = "Could not create board";
			console.error(err);
		});
	}

	const updateThisBoard = () => {
		if (props.currentBoardValue == undefined) {
			return;
		}

		BoardAPI.UpdateBoard(props.currentBoardValue.code, boardDTO.value).then((res: AxiosResponse<ApiResponse<BoardDTO>>) => {
			emits("created", res.data.data!.id);
		}).catch((err: AxiosError) => {
			error.value = `Failed to update ${props.currentBoardValue!.code}`;
		});
	}

	const setMimePreset = (presetName: string) => {
		switch (presetName) {
		case "image+video":
			boardDTO.value.config.mime_types_allowed = [];
			addMimeTypes(["image/jpeg", "image/png", "image/gif", "video/webm", "video/mp4"]);
			break;
		case "video+gif":
			boardDTO.value.config.mime_types_allowed = [];
			addMimeTypes(["video/webm", "video/mp4", "image/gif"]);
			break;
		case "image":
			boardDTO.value.config.mime_types_allowed = [];
			addMimeTypes(["image/jpeg", "image/png", "image/gif"]);
			break;
		case "video":
			boardDTO.value.config.mime_types_allowed = [];
			addMimeTypes(["video/webm", "video/mp4"]);
			break;
		}
	}

	const removeMimeTypeAtIndex = (index: number) => {
		boardDTO.value.config.mime_types_allowed.splice(index, 1);
	}

	const addMimeType = (mimeType: string) => {
		if (props.allowedMimeTypes.indexOf(mimeType) == -1) {
			return;
		}
		if (boardDTO.value.config.mime_types_allowed.indexOf(mimeType) != -1) {
			return;
		}
		boardDTO.value.config.mime_types_allowed.push(mimeType);
	}

	const addMimeTypes = (mimeTypes: string[]) => {
		for (let mimeType of mimeTypes) addMimeType(mimeType);
	}

	const onClickAddMimeType = () => {
		addMimeType(addMimeTypeInputVal.value);
		addMimeTypeInputVal.value = "";
	}
</script>

<template>
	<form @submit.prevent="onSubmit">
		<label for="code"><abbr title="Shorthand name for the board, used in URLs and when cross-referencing">Code</abbr>: </label>
		<input id="code" type=text placeholder="/a/, /b/, ..." required v-model="boardDTO.code"/><br/>

		<label for="name">Name: </label>
		<input id="name" type=text placeholder="('Anime and Manga')" required v-model="boardDTO.name"/><br/>

		<label for="description">Description: </label>
		<textarea id="description" placeholder="Description" v-model="boardDTO.description"/><br/>

		<label for="max-file-size">Max file size: </label>
		<input id="max-file-size" type=text placeholder="2 MB, 256KB, 1048576" required v-model="maxFileSizeInputVal"/><br/>

		<label for="reply-files-allowed">Replies can post files: </label>
		<input id="reply-files-allowed" type=checkbox v-model="boardDTO.config.reply_files_allowed" /><br/>

		<span class="input-group">
			<label>Allowed MIME types for files: </label>
			<i v-if="boardDTO.config.mime_types_allowed.length == 0">None</i>
			<ul>
				<li v-for="mimetype, i of boardDTO.config.mime_types_allowed">
					{{ mimetype }} [<a href="#" @click.prevent="removeMimeTypeAtIndex(i)">&#10006;</a>]
				</li>
			</ul>

			<button type="button" @click="onClickAddMimeType">+</button>
			<input id="allowed-mime-types" v-model="addMimeTypeInputVal" list="allowed-mime-types-datalist" placeholder="image/png (autocomplete)" />
			<datalist id="allowed-mime-types-datalist">
				<option
					v-for="mime of allowedMimeTypes.filter(x => !boardDTO.config.mime_types_allowed.includes(x))"
					:value="mime" />
			</datalist>

			<br/><br/>
			<label>Presets: </label>
			<button type="button" @click="setMimePreset('image+video')">Image and Video</button>
			<button type="button" @click="setMimePreset('video+gif')">Video and GIF</button>
			<button type="button" @click="setMimePreset('image')">Image only</button>
			<button type="button" @click="setMimePreset('video')">Video only</button>
		</span><br/>

		<label for="spoilers-enabled">Spoilers enabled: </label>
		<input id="spoilers-enabled" type=checkbox v-model="boardDTO.config.allow_spoilers" /><br/>

		<div v-if="boardDTO.config.allow_spoilers">
			<label for="spoiler-image">Spoilers image: </label>

			<div>
				<div v-for="fname of spoiler_images">
					<input type="radio" :id="`spoiler-${fname}`" :value="fname" v-model="boardDTO.config.spoiler_image" />
					<label :for="`spoiler-${fname}`">
						<img :src="CdnAPI.GetSpoilerURI(fname)" />
					</label>
				</div>
			</div>
		</div>

		<label for="locked">Locked: </label>
		<input id="locked" type=checkbox v-model="boardDTO.config.locked" /><br/>

		<label for="hidden">Hidden: </label>
		<input id="hidden" type=checkbox v-model="boardDTO.config.hidden" /><br/>

		<label for="max-threads">Maximum threads: </label>
		<input id="max-threads" type=number :min="1" :max="1000" v-model="boardDTO.config.max_threads" /><br/>

		<label for="bump-limit">Bump limit: </label>
		<input id="bump-limit" type=number :min="1" :max="1000" v-model="boardDTO.config.bump_limit" /><br/>

		<label for="image-limit">Image limit: </label>
		<input id="image-limit" type=number :min="1" :max="1000" v-model="boardDTO.config.image_limit" /><br/>

		<label for="flags-enabled">Flags enabled: </label>
		<input id="flags-enabled" type=checkbox v-model="boardDTO.config.flags_enabled" /><br/>

		<label for="ids-enabled">IDs enabled: </label>
		<input id="ids-enabled" type=checkbox v-model="boardDTO.config.ids_enabled" /><br/>

		<label for="math-enabled">Math typesetting enabled: </label>
		<input id="math-enabled" type=checkbox v-model="boardDTO.config.math_enabled" /><br/>

		<label for="code-enabled">Code blocks enabled: </label>
		<input id="code-enabled" type=checkbox v-model="boardDTO.config.code_enabled" /><br/>

		<br/>
		<button type=submit>
			<template v-if="isFormForNewBoard()">Create new Board</template>
			<template v-else>Update /{{props.currentBoardValue!.code}}/</template>
		</button>

		<template v-if="error">
			<div/>
			<span class="error">{{error}}</span>
		</template>
	</form>
</template>

<style scoped>
	.error {
		color: var(--user-error-color);
	}

	.input-group {
		border: 1px solid black;
		padding: 0.5em;
		display: inline-block;
		margin: 2px 0px;
	}
</style>
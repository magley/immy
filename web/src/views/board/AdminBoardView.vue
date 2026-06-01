<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { BoardAPI, type BoardDTO, type CreateBoardDTO, type UpdateBoardDTO } from "@/api/board.api.ts";
	import type { AxiosError, AxiosResponse } from 'axios';
	import type { ApiResponse } from '@/api/http';
	import { GetElipsisString, IsAlphaNumeric, StripSlashes } from '@/util/various.util';
	import { GetFileSizeByteFromString, GetFileSizeByteString } from '@/util/file.util';
	import { MetaAPI } from '@/api/meta.api';

	const boards = ref<BoardDTO[]>([]);
	const createBoardDTO = ref<CreateBoardDTO>({
		name: '',
		code: '',
		description: null,
		config: {
			locked: false,
			hidden: false,
			max_file_size: 0,
			reply_files_allowed: true,
			mime_types_allowed: [],
			bump_limit: 250,
			image_limit: 150,
			flags_enabled: false,
			ids_enabled: false
		},
	});
	const addMimeTypeInputVal = ref<string>("");
	const maxFileSizeInputVal = ref<string>("2 MB");
	const createBoardError = ref<string | undefined>(undefined);
	const allowedMimeTypes = ref<string[]>([]);
	
	onMounted(() => {
		get_boards();

		MetaAPI.GetMimeTypes().then((res: AxiosResponse<ApiResponse<string[]>>) => {
			allowedMimeTypes.value = res.data.data!;
			setMimePreset("image+video");
		}).catch((err) => {
			console.error(err);
		});
	});
	
	const get_boards = () => {
		BoardAPI.ListBoards().then((res: AxiosResponse<ApiResponse<BoardDTO[]>>) => {
			boards.value = res.data.data!;
		}).catch((err) => {
			console.error(err);
		});
	}
	
	const onSubmitCreateBoard = () => {
		createBoardDTO.value.code = StripSlashes(createBoardDTO.value.code.trim());
		createBoardDTO.value.name = createBoardDTO.value.name.trim();
		createBoardDTO.value.description = createBoardDTO.value.description?.trim() ?? null;
		createBoardDTO.value.config.max_file_size = GetFileSizeByteFromString(maxFileSizeInputVal.value);

		createBoardError.value = "";

		if (!IsAlphaNumeric(createBoardDTO.value.code)) {
			createBoardError.value = "Board code must consist of only alpha-numeric characters";
			return;
		}
		if (isNaN(createBoardDTO.value.config.max_file_size)) {
			createBoardError.value = "Invalid maximum file size";
			return;
		}
		if (createBoardDTO.value.config.mime_types_allowed.length == 0) {
			createBoardError.value = "No MIME types specified (users won't be able to attach files)";
			return;
		}

		BoardAPI.CreateBoard(createBoardDTO.value).then((res: AxiosResponse<ApiResponse<BoardDTO>>) => {
			get_boards();
		}).catch((err: AxiosError) => {
			createBoardError.value = "Could not create board";
			console.error(err);
		});
	}
	
	const onSubmitChangesToBoard = (idx: number) => {
		// Enter special view...
	}
	
	const onDeleteBoard = (idx: number) => {
		const board: BoardDTO = boards.value[idx]!;
		
		BoardAPI.DeleteBoard(board.code).then((res: AxiosResponse<ApiResponse<number>>) => {
			get_boards();
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const setMimePreset = (presetName: string) => {
		switch (presetName) {
		case "image+video":
			createBoardDTO.value.config.mime_types_allowed = [];
			addMimeTypes(["image/jpeg", "image/png", "image/gif", "video/webm", "video/mp4"]);
			break;
		case "video+gif":
			createBoardDTO.value.config.mime_types_allowed = [];
			addMimeTypes(["video/webm", "video/mp4", "image/gif"]);
			break;
		case "image":
			createBoardDTO.value.config.mime_types_allowed = [];
			addMimeTypes(["image/jpeg", "image/png", "image/gif"]);
			break;
		case "video":
			createBoardDTO.value.config.mime_types_allowed = [];
			addMimeTypes(["video/webm", "video/mp4"]);
			break;
		}
	}

	const removeMimeTypeAtIndex = (index: number) => {
		createBoardDTO.value.config.mime_types_allowed.splice(index, 1);
	}

	const addMimeType = (mimeType: string) => {
		if (allowedMimeTypes.value.indexOf(mimeType) == -1) {
			return;
		}
		if (createBoardDTO.value.config.mime_types_allowed.indexOf(mimeType) != -1) {
			return;
		}
		createBoardDTO.value.config.mime_types_allowed.push(mimeType);
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
	<h1>Boards</h1>
	
	<form @submit.prevent="onSubmitCreateBoard">
		<h2>Create new board</h2>
		<label for="code"><abbr title="Shorthand name for the board, used in URLs and when cross-referencing">Code</abbr>: </label>
		<input id="code" type=text placeholder="/a/, /b/, ..." required v-model="createBoardDTO.code"/><br/>

		<label for="name">Name: </label>
		<input id="name" type=text placeholder="('Anime and Manga')" required v-model="createBoardDTO.name"/><br/>

		<label for="description">Description: </label>
		<textarea id="description" placeholder="Description" v-model="createBoardDTO.description"/><br/>

		<label for="max-file-size">Max file size: </label>
		<input id="max-file-size" type=text placeholder="2 MB, 256KB, 1048576" required v-model="maxFileSizeInputVal"/><br/>

		<label for="reply-files-allowed">Replies can post files: </label>
		<input id="reply-files-allowed" type=checkbox v-model="createBoardDTO.config.reply_files_allowed" /><br/>

		<span class="input-group">
			<label>Allowed MIME types for files: </label>
			<i v-if="createBoardDTO.config.mime_types_allowed.length == 0">None</i>
			<ul>
				<li v-for="mimetype, i of createBoardDTO.config.mime_types_allowed">
					{{ mimetype }} [<a href="#" @click.prevent="removeMimeTypeAtIndex(i)">&#10006;</a>]
				</li>
			</ul>

			<button type="button" @click="onClickAddMimeType">+</button>
			<input id="allowed-mime-types" v-model="addMimeTypeInputVal" list="allowed-mime-types-datalist" placeholder="image/png (autocomplete)" />
			<datalist id="allowed-mime-types-datalist">
				<option
					v-for="mime of allowedMimeTypes.filter(x => !createBoardDTO.config.mime_types_allowed.includes(x))"
					:value="mime" />
			</datalist>

			<br/><br/>
			<label>Presets: </label>
			<button type="button" @click="setMimePreset('image+video')">Image and Video</button>
			<button type="button" @click="setMimePreset('video+gif')">Video and GIF</button>
			<button type="button" @click="setMimePreset('image')">Image only</button>
			<button type="button" @click="setMimePreset('video')">Video only</button>
		</span><br/>

		<label for="bump-limit">Bump limit: </label>
		<input id="bump-limit" type=number :min="1" :max="1000" v-model="createBoardDTO.config.bump_limit" /><br/>

		<label for="image-limit">Image limit: </label>
		<input id="image-limit" type=number :min="1" :max="1000" v-model="createBoardDTO.config.image_limit" /><br/>

		<label for="flags-enabled">Flags enabled: </label>
		<input id="flags-enabled" type=checkbox v-model="createBoardDTO.config.flags_enabled" /><br/>

		<label for="ids-enabled">IDs enabled: </label>
		<input id="ids-enabled" type=checkbox v-model="createBoardDTO.config.ids_enabled" /><br/>

		<br/>
		<button type=submit>Create</button>
		
		<template v-if="createBoardError">
			<div/>
			<span class="error">{{createBoardError}}</span>
		</template>
	</form>
	
	<br />

	<h2>Board list</h2>

	<table>
		<tbody>
			<tr>
				<th>ID</th>
				<th>Code</th>
				<th>Name</th>
				<th>Desc</th>

				<th>Locked</th>
				<th>Hidden</th>
				<th>Max File</th>
				<th><abbr title="Can non OPs post files in replies?">Reply Files</abbr></th>
				<th>MIME Types</th>
				<th>Bump Limit</th>
				<th>Image Limit</th>
				<th>Flags</th>
				<th>IDs</th>

				<th><b>Update</b></th>
				<th><b>Delete</b></th>
			</tr>
			<tr v-for="board, i in boards">
				<td>{{board.id}}</td>
				<td><RouterLink :to="`/${board.code}`">/{{board.code}}/</RouterLink></td>
				<td>{{board.name}}</td>
				<td>{{GetElipsisString(board.description)}}</td>
				<td>{{board.config.locked}}</td>
				<td>{{board.config.hidden}}</td>

				<td>{{GetFileSizeByteString(board.config.max_file_size)}}</td>
				<td>{{board.config.reply_files_allowed}}</td>
				<td><span v-for="mime, i of board.config.mime_types_allowed">{{mime}}, </span></td>
				<td>{{board.config.bump_limit}}</td>
				<td>{{board.config.image_limit}}</td>
				<td>{{board.config.flags_enabled}}</td>
				<td>{{board.config.ids_enabled}}</td>

				<td><button @click="onSubmitChangesToBoard(i)">Update...</button></td>
				<td><button @click="onDeleteBoard(i)">Delete</button></td>
			</tr>
		</tbody>
	</table>
</template>

<style scoped>
	table {
		border: 1px solid black;
	}
	
	th {
		background-color: var(--background-color-darker);
		padding: 0em 1em;
	}
	
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
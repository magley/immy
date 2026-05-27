<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { BoardAPI, type BoardDTO, type CreateBoardDTO, type UpdateBoardDTO } from "@/api/board.api.ts";
	import type { AxiosError, AxiosResponse } from 'axios';
	import type { ApiResponse } from '@/api/http';
import { GetElipsisString, IsAlphaNumeric, StripSlashes } from '@/util/various.util';
import { GetFileSizeByteString } from '@/util/file.util';

	const boards = ref<BoardDTO[]>([]);
	const createBoardDTO = ref<CreateBoardDTO>({
		name: '',
		code: '',
		description: null,
		config: {
			locked: false,
			hidden: false,
			max_file_size: 0,
			reply_files_allowed: false,
			mime_types_allowed: [],
			bump_limit: 0,
			image_limit: 0,
			flags_enabled: false,
			ids_enabled: false
		},
	});
	const createBoardError = ref<string | undefined>(undefined);
	
	onMounted(() => {
		get_boards();
	});
	
	const get_boards = () => {
		BoardAPI.ListBoards().then((res: AxiosResponse<ApiResponse<BoardDTO[]>>) => {
			boards.value = res.data.data!;
		}).catch((err) => {
			console.error(err);
		});
	}
	
	const onSubmitCreateBoard = () => {
		createBoardError.value = "";
		
		createBoardDTO.value.code = StripSlashes(createBoardDTO.value.code);

		if (!IsAlphaNumeric(createBoardDTO.value.code)) {
			createBoardError.value = "Board code must consist of only alpha-numeric characters";
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
</script>

<template>
	<h1>Boards</h1>
	
	<form @submit.prevent="onSubmitCreateBoard">
		<input type=text placeholder="Code (/a/, /b/, ...)" required v-model="createBoardDTO.code"/><br/>
		<input type=text placeholder="Name ('Anime and Manga')" required v-model="createBoardDTO.name"/><br/>
		<textarea placeholder="Description" v-model="createBoardDTO.description"/><br/>

		<label for="max-file-size">Max File Size: </label>
		<input id="max-file-size" type=text placeholder="2 MB" required v-model="createBoardDTO.name"/><br/>
		<br/>
		<button type=submit>Create board</button>
		
		<template v-if="createBoardError">
			<div/>
			<span class="error">{{createBoardError}}</span>
		</template>
	</form>
	
	<br />

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
				<td>/{{board.code}}/</td>
				<td>{{board.name}}</td>
				<td>{{GetElipsisString(board.description)}}</td>
				<td>{{board.config.locked}}</td>
				<td>{{board.config.hidden}}</td>

				<td>{{GetFileSizeByteString(board.config.max_file_size)}}</td>
				<td>{{board.config.reply_files_allowed}}</td>
				<td>{{board.config.mime_types_allowed}}</td>
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
		background: lightgray;
		padding: 0em 1em;
	}
	
	.error {
		color: red;
	}
</style>
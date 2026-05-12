<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { BoardAPI, BoardDTO, CreateBoardDTO, UpdateBoardDTO } from "@/api/board.api.ts";

	const boards = ref<BoardDTO[]>([]);
	const createBoardDTO = ref<CreateBoardDTO>({});
	const createBoardError = ref<string>(undefined);
	
	onMounted(() => {
		get_boards();
	});
	
	const get_boards = () => {
		BoardAPI.ListBoards().then((res: AxiosResponse<ApiResponse<BoardDTO[]>>) => {
			boards.value = res.data.data;
		}).catch((err) => {
			console.error(err);
		});
	}
	
	const onSubmitCreateBoard = () => {
		createBoardError.value = "";
		
		createBoardDTO.value.code = stripSlashes(createBoardDTO.value.code);
				
		if (!isAlphaNumeric(createBoardDTO.value.code)) {
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
	
	const onSubmitChangesToBoard = (idx: int) => {
		// Enter special view...
	}
	
	const onDeleteBoard = (idx: int) => {
		const board: BoardDTO = boards.value[idx];
		
		BoardAPI.DeleteBoard(board.code).then((res: AxiosResponse<ApiResponse<number>>) => {
			get_boards();
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const getElipsisString = (input: string, lenIncludingDots: number = 10 + 3) => {
		if (!input) return input;
		
		if (input.length > lenIncludingDots) {
			return `${input.substring(0, lenIncludingDots)}...`;
		} else {
			return input;
		}
	}
	
	const isAlphaNumeric = (str: string) =>  {
		var code, i, len;

		for (i = 0, len = str.length; i < len; i++) {
			code = str.charCodeAt(i);
			if (!(code > 47 && code < 58) && !(code > 64 && code < 91) && !(code > 96 && code < 123)) {
				return false;
			}
		}
		return true;
	}
	
	const stripSlashes = (str: string) => {
		if (str.startsWith("/") && str.endsWith("/")) {
			return str.slice(1, -1);
		}
		return str;
	}
</script>

<template>
	<h1>Boards</h1>
	
	<form @submit.prevent="onSubmitCreateBoard">
		<input type=text placeholder="Code (/a/, /b/, ...)" required v-model="createBoardDTO.code"/><br/>
		<input type=text placeholder="Name ('Anime and Manga')" required v-model="createBoardDTO.name"/><br/>
		<textarea placeholder="Description" v-model="createBoardDTO.description"/><br/>
		<br/>
		<button type=submit>Create board</button>
		
		<template v-if="createBoardError">
			<div/>
			<span class="error">{{createBoardError}}</span>
		</template>
	</form>
	
	<br />

	<table>
		<tr>
			<th>ID</th>
			<th>Code</th>
			<th>Name</th>
			<th>Description</th>
			<th>Locked</th>
			<th>Hidden</th>
			<th>Update</th>
			<th>Delete</th>
		</tr>
		<tr v-for="board, i in boards">
			<td>{{board.id}}</td>
			<td>/{{board.code}}/</td>
			<td>{{board.name}}</td>
			<td>{{getElipsisString(board.description)}}</td>
			<td>{{board.locked}}</td>
			<td>{{board.hidden}}</td>
			<td>
				<button @click="onSubmitChangesToBoard(i)">Update...</button>
			</td>
			<td>
				<button @click="onDeleteBoard(i)">Delete</button>
			</td>
		</tr>
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
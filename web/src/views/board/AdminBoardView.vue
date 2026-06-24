<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { BoardAPI, type BoardDTO } from "@/api/board.api.ts";
	import type { AxiosError, AxiosResponse } from 'axios';
	import type { ApiResponse } from '@/api/http';
	import { GetElipsisString } from '@/util/various.util';
	import { GetFileSizeByteString } from '@/util/file.util';
	import { MetaAPI } from '@/api/meta.api';
	import BoardUpdate from '@/components/board/BoardUpdate.vue';
	import { UserAPI, UserRole } from '@/api/user.api';
	import { useRouter } from 'vue-router';

	const router = useRouter();

	const boards = ref<BoardDTO[]>([]);
	const allowedMimeTypes = ref<string[]>([]);
	const boardEditShown = ref<Record<string, boolean>>({}); // board code => true/false
	
	onMounted(() => {
		UserAPI.AuthorizeUser({required_roles: [UserRole.Admin]}).then(() => {
			get_boards();
			get_allowed_mime_types();
		}).catch((err: AxiosError) => {
			console.error(err);
			router.push("/login");
		});
	});

	const get_allowed_mime_types = () => {
		MetaAPI.GetMimeTypes().then((res: AxiosResponse<ApiResponse<string[]>>) => {
			allowedMimeTypes.value = res.data.data!;
		}).catch((err) => {
			console.error(err);
		});
	}
	
	const get_boards = () => {
		BoardAPI.ListBoards().then((res: AxiosResponse<ApiResponse<BoardDTO[]>>) => {
			boards.value = res.data.data!;
		}).catch((err) => {
			console.error(err);
		});
	}

	const onSubmitChangesToBoard = (board_code: string) => {
		boardEditShown.value[board_code] = !(boardEditShown.value[board_code] ?? false);
	}
	
	const onDeleteBoard = (idx: number) => {
		const board: BoardDTO = boards.value[idx]!;
		
		BoardAPI.DeleteBoard(board.code).then((res: AxiosResponse<ApiResponse<number>>) => {
			get_boards();
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const onBoardCreated = (board_id: number) => {
		get_boards();
	}

	const onBoardUpdated = (board_id: number) => {
		BoardAPI.GetBoardById(board_id).then((res) => {
			// Update the ref of the just-updated board.
			for (let i = 0; i < boards.value.length; i++) {
				if (boards.value[i]!.id == board_id) {
					boards.value[i] = res.data.data!;
					boardEditShown.value[res.data.data!.code] = false;
					break;
				}
			}
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}
</script>

<template>
	<h1>Boards</h1>

	<div id="new-board-container">
		<h3>Create new Board</h3>
		<BoardUpdate
		:allowed-mime-types="allowedMimeTypes"
		:current-board-value="undefined"
		@created="onBoardCreated"
		/>
	</div>
	
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
				<th>Max Threads</th>
				<th>Bump Limit</th>
				<th>Image Limit</th>
				<th>Spoilers</th>
				<th>Flags</th>
				<th>IDs</th>
				<th>Math</th>
				<th>Code</th>

				<th><b>Update</b></th>
				<th><b>Delete</b></th>
			</tr>
			<template v-for="board, i in boards">
				<tr>
					<td>{{board.id}}</td>
					<td><RouterLink :to="`/${board.code}`">/{{board.code}}/</RouterLink></td>
					<td>{{board.name}}</td>
					<td>{{GetElipsisString(board.description)}}</td>
					<td><img src="/icons/lock.png" v-if="board.config.locked" /></td>
					<td><img src="/icons/invisible.png" v-if="board.config.hidden" /></td>

					<td>{{GetFileSizeByteString(board.config.max_file_size)}}</td>
					<td>{{board.config.reply_files_allowed}}</td>
					<td><span v-for="mime, i of board.config.mime_types_allowed"></span></td>
					<td>{{board.config.max_threads}}</td>
					<td>{{board.config.bump_limit}}</td>
					<td>{{board.config.image_limit}}</td>
					<td>{{board.config.allow_spoilers}}</td>
					<td>{{board.config.flags_enabled}}</td>
					<td>{{board.config.ids_enabled}}</td>
					<td>{{board.config.math_enabled}}</td>
					<td>{{board.config.code_enabled}}</td>

					<td><button @click="onSubmitChangesToBoard(board.code)">Update...</button></td>
					<td><button @click="onDeleteBoard(i)">Delete</button></td>
				</tr>
				<tr v-if="boardEditShown[board.code] ?? false">
					<td colspan="17">
						<BoardUpdate
						:current-board-value="board"
						:allowed-mime-types="allowedMimeTypes"
						@created="onBoardUpdated"
						class="update-board-form"
						/>
					</td>
				</tr>
			</template>
		</tbody>
	</table>
</template>

<style scoped>
	h1, h2, h3 {
		color: var(--banner-title-color);
	}

	table {
		border: 1px solid black;

		th {
			background-color: var(--background-color-accent);
			padding: 0.5em 0.5em;
		}

		td {
			padding: 0.5em 0;
			img {
				display: block;
				margin: auto;
			}
		}
	}
	


	#new-board-container {
		text-align: center;

		padding: 1em;
		border: 1px solid black;
		background-color: var(--background-color-darker);
		display: block;
		margin: auto;
		width: 40%;
	}

	.update-board-form {
		padding: 1em;
		border: 1px solid black;
		background-color: var(--background-color-darker);
		display: block;
		margin: auto;
		width: 50%;
	}
</style>
<script setup lang="ts">
	import { BoardAPI, type BoardDTO } from '@/api/board.api';
	import type { AxiosError } from 'axios';
	import { onMounted, ref } from 'vue';


	const boards = ref<BoardDTO[]>([]);

	onMounted(() => {
		loadBoards();
	});

	const loadBoards = () => {
		// TODO: No limit here...
		BoardAPI.ListBoards(0, 1000).then((res) => {
			boards.value = res.data.data!;
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}
</script>

<template>
	<div id="title">
		<h1>
			<img src="/icons/immy-icon-48.png" />
			ImmyChan
		</h1>
		<p>
			<img src="/icons/immy-icon-16.png" />
			ImmyChan is an imageboard powered by the
			<a href="http://github.com/magley/immy" target="_blank" class="external">
				<img src="/icons/immy-icon-16.png" /> Immy imageboard engine
			</a>.
		</p>
	</div>

	<div class="container">
		<div class="header">
			<h2>Boards</h2>
		</div>
		<div class="body">
			<table>
				<tr v-for="board of boards">
					<td>
						<RouterLink :to="`/${board.code}`">
							/{{ board.code }}/ - {{ board.name }}
						</RouterLink>
						<span v-if="board.description != ''"><br/>{{ board.description }}</span>
					</td>
					<td>
						[<RouterLink :to="`/${board.code}/catalog`">Catalog</RouterLink>]
					</td>
				</tr>
			</table>
		</div>
	</div>

	<div class="container">
		<div class="header">
			<h2>Stats</h2>
		</div>
		<div class="body">

		</div>	</div>
</template>

<style scoped>
	#title {
		h1 {
			font-size: 4em;
			img {
				image-rendering: pixelated;
				vertical-align: middle;
			}
		}
	}

	.container {
		border: 1px solid var(--banner-title-color);
		border-radius: 15px;
		width: 50%;
		margin: 2em auto;

		background-color: white;

		.header {
			background-color: var(--highlighted-post-background-color);
			border-radius: 15px 15px 0 0;

			h2 {
				color: var(--banner-title-color);
				text-align: center;
				margin: 0;
				padding: 0.25em;
			}
		}

		.body {
			background-color: white;
			border-radius: 15px;
			padding: 1em;

			table {
				width: 100%;

				tr {
					td {
						margin: 0.5em 0;
					}
				}
			}
		}


	}
</style>
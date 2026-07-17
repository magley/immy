<script setup lang="ts">
	import { BoardAPI, type BoardDTO, type BoardStatisticsDTO } from '@/api/board.api';
	import { MetaAPI } from '@/api/meta.api';
import { GetFileSizeByteString } from '@/util/file.util';
	import type { AxiosError } from 'axios';
	import { onMounted, ref } from 'vue';

	const boards = ref<BoardDTO[]>([]);
	const stats_raw = ref<BoardStatisticsDTO[]>([]);

	const board_stats = ref<Record<string, BoardStatisticsDTO>>({});

	onMounted(() => {
		loadBoards();
		loadStats();
	});

	const loadBoards = () => {
		BoardAPI.GetAllBoards().then((res) => {
			boards.value = res.data.data!;
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}

	const loadStats = () => {
		MetaAPI.GetStatistics().then((res) => {
			stats_raw.value = res.data.data!;

			board_stats.value = {};
			for (let s of stats_raw.value) {
				board_stats.value[s.code] = s;
			}
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
					<td>
						{{ board_stats[board.code]?.thread_count_alive }} active threads
					</td>					<td>
						{{ board_stats[board.code]?.thread_count }} threads ever
					</td>
					<td>
						{{ board_stats[board.code]?.post_count }} posts
					</td>
				</tr>
			</table>
		</div>
	</div>

	<div class="container">
		<div class="header">
			<h2>Stats</h2>
		</div>
		<div class="body flex">
			<span>
				<b>Active threads:</b>
				{{ stats_raw.map((s) => s.thread_count_alive).reduce((a, b) => a + b, 0) }}
			</span>
			<span>
				<b>Total threads:</b>
				{{ stats_raw.map((s) => s.thread_count).reduce((a, b) => a + b, 0) }}
			</span>
			<span>
				<b>Total posts:</b>
				{{ stats_raw.map((s) => s.post_count).reduce((a, b) => a + b, 0) }}
			</span>

			<span>
				<b>Uploaded data:</b>
				{{ GetFileSizeByteString(stats_raw.map((s) => s.bytes_uploaded).reduce((a, b) => a + b, 0)) }}
			</span>
		</div>
	</div>

	<div class="container">
		<div class="body">
			<span>
				<RouterLink to="/blog">Blog</RouterLink>
			</span>
			•
			<span>
				<RouterLink to="/bans">Bans</RouterLink>
			</span>
		</div>
	</div>
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

	.flex {
		display: flex;

		span {
			flex: 1;
		}
	}
</style>
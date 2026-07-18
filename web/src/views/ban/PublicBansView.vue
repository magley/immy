<script setup lang="ts">
	import { BanAPI, type BanDTO } from '@/api/ban.api';
	import PaginatorComponent from '@/components/PaginatorComponent.vue';
	import { GetPostTimeReadable } from '@/model/post/post.model';
	import { Paginator } from '@/util/pagination.util';
	import { GetTimeDifferenceBasic } from '@/util/various.util';
	import { onMounted, reactive, ref } from 'vue';

	const error = ref<string | undefined>(undefined);
	const bans = ref<BanDTO[]>([]);
	const paginator = reactive<Paginator<BanDTO[]>>(new Paginator(BanAPI.ListBans));

	onMounted(() => {
		getBans();
	});

	const getBans = () => {
		paginator.getItems()
			.then((res) => bans.value = res.data.data!)
			.catch((_) => error.value = "Could not fetch bans!");
	}

	const gotoPage = (p: number) => {
		paginator.page = p;
		getBans();
	}
</script>

<template>
	<h1>Bans</h1>

	<div class="error" v-if="error">{{ error }}</div>
	<div v-if="paginator.loading">Loading...</div>
	<div>
		<div class="center nav">
			<PaginatorComponent :paginator="paginator" @goto-page="gotoPage" empty-message="There are no bans at the moment" />
			<br/>
		</div>

		<table>
			<tbody>
				<tr>
					<th>Board</th>
					<th>Duration</th>
					<th>Reason</th>
				</tr>
				<tr v-for="ban of bans">
					<td class="center">
						{{ ban.board_id == null ? "global" : `/${ban.board_code}/` }}
					</td>
					<td class="center">
						<template v-if="ban.warning">
							Warning
						</template>
						<template v-else-if="ban.expires_at">
							<span :title="`Expires on ${GetPostTimeReadable(ban.expires_at)}`">
								{{ GetTimeDifferenceBasic(new Date(Date.parse(ban.created_at)), new Date(Date.parse(ban.expires_at))) }}
							</span>
						</template>
						<template v-else>
							<b class="bad">Permanent</b>
						</template>
					</td>
					<td>{{ ban.reason }}</td>
				</tr>
			</tbody>
		</table>

		<div class="center nav">
			<PaginatorComponent :paginator="paginator" @goto-page="gotoPage" empty-message="There are no bans at the moment" />
		</div>
	</div>
</template>

<style scoped>
	.error {
		color: var(--user-error-color);
	}

	.bad {
		color: var(--user-error-color);
	}

	.center {
		text-align: center;
	}

	.nav {
		margin-top: 1em;
	}

	tt {
		font-size: 12pt;
	}

	h1 {
		color: var(--banner-title-color);
		text-align: center;
	}

	table {
		width: 50%;
		margin: auto;
		border: 1px solid black;
		background: white;

		th {
			background-color: var(--background-color-accent);
			padding: 0.5em 0.5em;
		}

		tr {
			margin: 0;
		}

		td {
			padding: 0.5em;
			border: 1px solid black;
			margin: 0;

			img {
				display: block;
				margin: auto;
			}
		}
	}
</style>
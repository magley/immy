<script setup lang="ts">
	import { BanAPI, type BanExtDTO } from '@/api/ban.api';
	import { UserAPI, UserRole } from '@/api/user.api';
	import { GetPostTimeReadable } from '@/model/post/post.model';
	import { GetTimeDifferenceBasic } from '@/util/various.util';
	import type { AxiosError } from 'axios';
	import { onMounted, ref } from 'vue';
	import { useRouter } from 'vue-router';

	const router = useRouter();
	const error = ref<string | undefined>(undefined);
	const loading = ref<boolean>(true);
	const bans = ref<BanExtDTO[]>([]);

	const perPage = 10;
	const page = ref<number>(1);
	const pagesTotal = ref<number>(1);
	const pagesNav = ref<number[]>([]);

	onMounted(() => {
		UserAPI.AuthorizeUser({required_roles: [UserRole.Admin, UserRole.Moderator]}).then(() => {
			getBans();
		}).catch((err: AxiosError) => {
			console.error(err);
			router.push("/login");
		});
	});

	const getBans = () => {
		if (page.value < 1) page.value = 1;
		if (page.value > pagesTotal.value) page.value = pagesTotal.value;

		loading.value = true;
		BanAPI.ListBansForAdmin((page.value - 1) * perPage, perPage).then((res) => {
			bans.value = res.data.data!;
			loading.value = false;

			const meta = res.data.meta!;
			page.value = meta.page;
			pagesTotal.value = meta.total_pages;
			pagesNav.value = [
				page.value - 4, page.value - 3, page.value - 2, page.value - 1,
				page.value - 0,
				page.value + 1, page.value + 2, page.value + 3, page.value + 4, page.value + 5,
			];
			pagesNav.value = pagesNav.value.filter((v) => v >= 1 && v <= meta.total_pages);
		}).catch((err: AxiosError) => {
			error.value = "Could not fetch bans";
		});
	}

	const gotoPage = (p: number) => {
		page.value = p;
		getBans();
	}
</script>

<template>
	<h1>Bans</h1>

	<div class="error" v-if="error">{{ error }}</div>
	<div v-if="loading">Loading...</div>
	<div>
		<table>
			<tbody>
				<tr>
					<th>Board</th>
					<th>Duration</th>
					<th>Reason</th>
				</tr>
				<tr v-for="ban, i of bans">
					<td class="center">
						{{ ban.board_code == null ? "global" : `/${ban.board_code}/` }}
					</td>
					<td class="center">
						<template v-if="ban.ban.warning">
							Warning
						</template>
						<template v-else-if="ban.ban.expires_at">
							<span :title="`Expires on ${GetPostTimeReadable(ban.ban.expires_at)}`">
								{{ GetTimeDifferenceBasic(new Date(Date.parse(ban.ban.created_at)), new Date(Date.parse(ban.ban.expires_at))) }}
							</span>
						</template>
						<template v-else>
							<b class="bad">Permanent</b>
						</template>
					</td>
					<td>{{ ban.ban.reason }}</td>
				</tr>
			</tbody>
		</table>

		<div class="center nav">
			[<a href="#" @click.prevent="gotoPage(1)">First</a>]&thinsp;
			[<a href="#" @click.prevent="gotoPage(page - 1)">Prev</a>]&thinsp;
			<span v-for="p, i of pagesNav">
				<template v-if="p == page">
					<span>{{ p }} </span>
				</template>
				<template v-else>
					<a href="#" @click.prevent="gotoPage(p)">{{ p }} </a>
				</template>
				<template v-if="i < pagesNav.length - 1">,</template>&thinsp;
			</span>
			[<a href="#" @click.prevent="gotoPage(page + 1)">Next</a>]&thinsp;
			[<a href="#" @click.prevent="gotoPage(pagesTotal)">Last</a>]&thinsp;
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
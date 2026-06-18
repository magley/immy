<script setup lang="ts">
	import { BanAPI, type BanDTO } from '@/api/ban.api';
	import { UserAPI, UserRole } from '@/api/user.api';
	import { GetPostTimeReadable } from '@/model/post/post.model';
	import type { AxiosError } from 'axios';
	import { onMounted, ref } from 'vue';
	import { useRouter } from 'vue-router';

	const router = useRouter();
	const error = ref<string | undefined>(undefined);
	const loading = ref<boolean>(true);
	const bans = ref<BanDTO[]>([]);

	onMounted(() => {
		// TODO: Authorize with mods too?
		UserAPI.AuthorizeUser({role: UserRole.Admin}).then(() => {
			getBans();
		}).catch((err: AxiosError) => {
			console.error(err);
			router.push("/login");
		});
	});

	const getBans = () => {
		loading.value = true;
		BanAPI.ListBansForAdmin().then((res) => {
			bans.value = res.data.data!;
			loading.value = false;
		}).catch((err: AxiosError) => {
			error.value = "Could not fetch bans";
		});
	}

	const banNumToIP = (num: number): string => {
		return [
			(num >>> 24) & 255,
			(num >>> 16) & 255,
			(num >>> 8) & 255,
			num & 255
		].join('.');
	}

	const isBanFinished = (ban: BanDTO) => {
		return isExpired(ban) || isWarningSeen(ban) || isRemoved(ban);
	}

	const isExpired = (ban: BanDTO) => {
		if (!ban.expires_at) return false;
		const expiresAt = Date.parse(ban.expires_at);
		if (expiresAt > new Date().getMilliseconds()) return true;
		return false;
	}

	const isWarningSeen = (ban: BanDTO): boolean => {
		return (ban.warning && ban.seen);
	}

	const isRemoved = (ban: BanDTO) => {
		return (ban.deleted_at != null);
	}
</script>

<template>
	<h1>Bans</h1>

	<div class="error" v-if="error">{{ error }}</div>
	<div v-if="loading">Loading...</div>
	<div v-else>
		<table>
			<tbody>
				<tr>
					<th>ID</th>
					<th>IP</th>
					<th>Issued</th>
					<th>Expires</th>
					<th>By</th>
					<th>Reason</th>
					<th></th>
				</tr>
				<tr v-for="ban, i of bans" :class="{finished: isBanFinished(ban)}">
					<td>{{ ban.id }}</td>
					<td>
						<template v-if="ban.ip_end == null">
							{{ banNumToIP(ban.ip_start) }}
						</template>
						<template v-else>
							{{ banNumToIP(ban.ip_start) }} - {{ banNumToIP(ban.ip_end) }}
						</template>
					</td>
					<td>{{ GetPostTimeReadable(ban.created_at) }}</td>
					<td>
						<template v-if="ban.warning">
							Warning
						</template>
						<template v-else-if="ban.expires_at">
							{{ GetPostTimeReadable(ban.expires_at) }}
						</template>
						<template v-else>
							<b class="bad">Permanent</b>
						</template>
					</td>
					<td>{{ ban.creator_id }}</td>
					<td>{{ ban.reason }}</td>
					<td>
						<template v-if="isExpired(ban)">
							<img src="/icons/clock.png" :title="`Expired at ${GetPostTimeReadable(ban.expires_at!)}.`" />
						</template>
						<template v-else-if="isWarningSeen(ban)">
							<img src="/icons/visible.png" :title="`User acknowledged this warning`" />
						</template>
						<template v-else-if="isRemoved(ban)">
							<img src="/icons/trash.png" :title="`Ban manually removed at ${GetPostTimeReadable(ban.deleted_at!)}`" />
						</template>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>

<style scoped>
	.error {
		color: var(--user-error-color);
	}

	.bad {
		color: var(--user-error-color);
	}

	h1 {
		color: var(--banner-title-color);
		text-align: center;
	}

	table {
		width: 80%;
		margin: auto;
		border: 1px solid black;

		th {
			background-color: var(--background-color-accent);
			padding: 0.5em 0.5em;
		}

		tr.finished > td {
			background-color: lightgray;
		}

		td {
			padding: 0.5em 0;
			img {
				display: block;
				margin: auto;
			}
		}
	}
</style>
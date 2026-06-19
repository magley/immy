<script setup lang="ts">
	import { BanAPI, type BanDTO, type BanExtDTO } from '@/api/ban.api';
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

	onMounted(() => {
		UserAPI.AuthorizeUser({required_roles: [UserRole.Admin, UserRole.Moderator]}).then(() => {
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

	const isBanFinished = (ban: BanExtDTO) => {
		return isExpired(ban) || isWarningSeen(ban) || isRemoved(ban);
	}

	const isExpired = (ban: BanExtDTO) => {
		if (!ban.ban.expires_at) return false;
		const expiresAt = Date.parse(ban.ban.expires_at);
		if (expiresAt < new Date().getTime()) return true;
		return false;
	}

	const isWarningSeen = (ban: BanExtDTO): boolean => {
		return (ban.ban.warning && ban.ban.seen);
	}

	const isRemoved = (ban: BanExtDTO) => {
		return (ban.ban.deleted_at != null);
	}

	const removeBan = (ban: BanExtDTO) => {
		BanAPI.DeleteBan(ban.ban.id).then((res) => {
			refreshBan(ban.ban.id);
		}).catch((err: AxiosError) => {
			error.value = `Failed to delete ban #${ban.ban.id}`;
			console.error(err);
		});
	}

	const refreshBan = (banID: number) => {
		BanAPI.GetBanForAdmin(banID).then((res) => {
			for (let i = 0; i < bans.value.length; i++) {
				if (bans.value[i]!.ban.id == banID) {
					bans.value[i] = res.data.data!
				}
			}
		}).catch((err: AxiosError) => {
			error.value = `Failed to refresh ban #${banID}`;
			console.error(err);
		})
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
					<th>Board</th>
					<th>Issued</th>
					<th>Duration</th>
					<th>By</th>
					<th>Reason</th>
					<th></th>
				</tr>
				<tr v-for="ban, i of bans" :class="{finished: isBanFinished(ban)}">
					<td>{{ ban.ban.id }}</td>
					<td>
						<tt>
							<template v-if="ban.ban.ip_end == null">
								{{ banNumToIP(ban.ban.ip_start) }}
							</template>
							<template v-else>
								{{ banNumToIP(ban.ban.ip_start) }} - {{ banNumToIP(ban.ban.ip_end) }}
							</template>
						</tt>
					</td>
					<td class="center">
						{{ ban.board_code == null ? "global" : `/${ban.board_code}/` }}
					</td>
					<td>{{ GetPostTimeReadable(ban.ban.created_at) }}</td>
					<td>
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
					<td class="center">{{ ban.creator_username }}</td>
					<td>{{ ban.ban.reason }}</td>
					<td>
						<template v-if="isExpired(ban)">
							<img src="/icons/clock.png" :title="`Expired at ${GetPostTimeReadable(ban.ban.expires_at!)}.`" />
						</template>
						<template v-else-if="isWarningSeen(ban)">
							<img src="/icons/visible.png" :title="`User acknowledged this warning`" />
						</template>
						<template v-else-if="isRemoved(ban)">
							<img src="/icons/trash.png" :title="`Ban manually removed at ${GetPostTimeReadable(ban.ban.deleted_at!)}`" />
						</template>
						<template v-else>
							<button @click="removeBan(ban)" :title="`Remove ban`"><img src="/icons/delete.png"/></button>
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

	.center {
		text-align: center;
	}

	tt {
		font-size: 12pt;
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
			background-color: var(--post-background-color);
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
<script setup lang="ts">
	import { BanAPI, type BanExtDTO } from '@/api/ban.api';
	import { BanAppealAPI, type BanAppealDTO, BanAppealStatus, type UpdateBanAppealDTO } from '@/api/ban_appeal.api';
	import { UserAPI, UserRole } from '@/api/user.api';
	import { GetPostTimeReadable } from '@/model/post/post.model';
	import { GetTimeDifferenceBasic } from '@/util/various.util';
	import type { AxiosError } from 'axios';
	import { onMounted, reactive, ref } from 'vue';
	import { useRoute, useRouter } from 'vue-router';
	import PaginatorComponent from '@/components/PaginatorComponent.vue';
	import { Paginator } from '@/util/pagination.util';

	const route = useRoute();
	const router = useRouter();
	const error = ref<string | undefined>(undefined);
	const bans = ref<BanExtDTO[]>([]);
	const paginator = reactive<Paginator<BanExtDTO[]>>(new Paginator(BanAPI.ListBansForAdmin));
	const banWhoseAppealIsOpen = ref<number | undefined>(undefined);
	const appeals = ref<Record<number, BanAppealDTO[]>>({});
	const appealError = ref<string | undefined>(undefined);
	const banAppealResponse = ref<BanAppealStatus>(BanAppealStatus.Rejected);

	onMounted(() => {
		UserAPI.AuthorizeUser({required_roles: [UserRole.Admin, UserRole.Moderator]}).then(() => {
			getBans();
		}).catch((err: AxiosError) => {
			console.error(err);
			router.push({ path: '/login', query: { redirect: route.fullPath } });
		});
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

	const onClickBanAppealShowButton = (ban: BanExtDTO) => {
		if (banWhoseAppealIsOpen.value == ban.ban.id) {
			banWhoseAppealIsOpen.value = undefined;
		} else {
			banWhoseAppealIsOpen.value = ban.ban.id;
			appealError.value = undefined;
			loadBanAppeals(ban.ban.id);
		}
	}

	const loadBanAppeals = (banID: number) => {
		appealError.value = undefined;
		BanAppealAPI.GetBanAppealsOfBan(banID).then((res) => {
			appeals.value[banID] = res.data.data!;
		}).catch((err: AxiosError) => {
			appealError.value = "Could not load ban appeals";
			console.error(err);
		});
	}

	const submitAppealResponse = (ban: BanExtDTO, appeal: BanAppealDTO) => {
		const dto: UpdateBanAppealDTO = {
			status: banAppealResponse.value
		};
		BanAppealAPI.UpdateBanAppeal(appeal.id, dto).then(() => {
			loadBanAppeals(ban.ban.id);
		}).catch((err: AxiosError) => {
			appealError.value = "Could not submit your response to the ban appeal";
			console.error(err);
		});
	}
</script>

<template>
	<h1>Bans</h1>

	<div class="error" v-if="error">{{ error }}</div>
	<div v-if="paginator.loading">Loading...</div>
	<div v-else>
		<div class="center nav">
			<PaginatorComponent :paginator="paginator" @goto-page="gotoPage" />
			<br/>
		</div>

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
					<th>Appeal</th>
				</tr>
				<template v-for="ban, i of bans">
					<tr :class="{finished: isBanFinished(ban)}">
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
						<td>
							<button @click="onClickBanAppealShowButton(ban)" v-if="!ban.ban.warning">Show/Hide</button>
						</td>
					</tr>
					<tr v-if="banWhoseAppealIsOpen == ban.ban.id">
						<td colspan="9">
							<h3>Appeals:</h3>
							<span v-if="appealError" class="error">{{ error }}</span>
							<div>
								<div v-for="appeal, i of appeals[ban.ban.id]" class="ban-appeal-message">
									<div>
										<span class="user">User</span>
										{{ GetPostTimeReadable(appeal.created_at) }}
										<br /><br />
										{{ appeal.message }}
									</div>
									<template v-if="i > 0 || appeal.status != BanAppealStatus.Pending || isBanFinished(ban)">
										<br /><br />
										<template v-if="appeal.status == BanAppealStatus.Pending">
											<i>This appeal was never handled.</i>
										</template>
										<template v-else>
											<i>Responsed with status</i> <tt><b>{{ appeal.status }}</b></tt>
										</template>
									</template>
									<template v-else>
										<br />
										<div class="ban-appeal-response-form">
											<input  id="appeal-approve" type="radio" :value="BanAppealStatus.Approved" v-model="banAppealResponse" />
											<label for="appeal-approve">Approve</label>
											<br />

											<input  id="appeal-reject" type="radio" :value="BanAppealStatus.Rejected" v-model="banAppealResponse" />
											<label for="appeal-reject">Reject</label>
											<br />

											<input  id="appeal-reject-final" type="radio" :value="BanAppealStatus.RejectedFinal" v-model="banAppealResponse" />
											<label for="appeal-reject-final">Reject (no more reappeals)</label>
											<br />
											<br/>
											<button @click="submitAppealResponse(ban, appeal)">Submit response</button>
										</div>
									</template>
								</div>
							</div>
						</td>
					</tr>
				</template>
			</tbody>
		</table>

		<div class="center nav">
			<PaginatorComponent :paginator="paginator" @goto-page="gotoPage" />
		</div>
	</div>
</template>

<style scoped>
	.error {
		color: var(--user-error-color);
	}

	.nav {
		margin-top: 1em;
	}

	.ban-appeal-message {
		background-color: var(--post-background-color);
		padding: 1em;
		margin: 1.2em;

		.user {
			color: var(--username-color);
			font-weight: bold;
		}

		.ban-appeal-response-form {
			display: inline-block;
			border: 1px solid black;
			padding: 0.8em;
			margin: 0.2em;
			background: var(--background-color);
		}
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
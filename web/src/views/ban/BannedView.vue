<script setup lang="ts">
	import { BanAPI, type BanDTO, type UpdateBanDTO } from '@/api/ban.api';
	import type { ApiResponse } from '@/api/http';
	import { GetPostTimeReadable } from '@/model/post/post.model';
	import type { AxiosError, AxiosResponse } from 'axios';
	import { onMounted, ref, warn } from 'vue';

	const bans = ref<BanDTO[]>([]);
	const error = ref<string | undefined>(undefined);
	const loading = ref<boolean>(true);

	onMounted(() => {
		getBans();
	});

	const getBans = () => {
		loading.value = true;
		error.value = undefined;
		BanAPI.GetMyBans().then((res: AxiosResponse<ApiResponse<BanDTO[]>>) => {
			bans.value = res.data.data!;
			markExpiredBansAsSeen();
		}).catch((err: AxiosError) => {
			error.value = "Could not fetch your bans";
			console.error(err);
		}).finally(() => {
			loading.value = false;
		});
	}

	const markExpiredBansAsSeen = () => {
		for (let ban of bans.value) {
			if (expired(ban)) {
				const dto: UpdateBanDTO = {
					seen: true,
				};
				BanAPI.UpdateBan(ban.id, dto);
			}
		}
	}

	const expired = (ban: BanDTO): boolean => {
		if (ban.warning) return true;
		if (ban.expires_at == null) return false;
		return Date.parse(ban.expires_at) < new Date().getTime();
	}
</script>

<template>
	<div class="container">
		<div class="header">
			<h2>
				Banned<template v-if="bans.length == 0">?</template>
				<template v-else>! ;-;</template>
			</h2>
		</div>
		<div class="body">
			<div v-if="loading">
				Checking to see if you are banned...
			</div>
			<div v-else>
				<div v-if="error">
					<b class="error">
						{{ error }}
					</b>
				</div>
				<div v-else>
					<div v-if="bans.length == 0">
						You are currently not banned <img src="/icons/immy-icon-16.png" />!
						<br/>
						Return <RouterLink to="/">home</RouterLink>.
					</div>
					<div v-else>
						<div v-for="ban, i of bans" class="ban-info">
							<p>
								You have <u v-if="i > 0">also</u> been
								<template v-if="ban.warning">
									<b>warned</b>
								</template>
								<template v-else>
									<b>banned</b> from
									<template v-if="ban.board_id == null">
										<b>all boards</b>
									</template>
									<template v-else>
										<b>{{ ban.board_id }}</b>
									</template>
								</template>
								for the following reason:
							</p>
							<div class="ban-reason">
								{{ ban.reason }}
							</div>
							<p>
								This
								<template v-if="ban.warning">warning</template><template v-else>ban</template>
								has been issued at
								<b>{{ GetPostTimeReadable(ban.created_at) }}</b>
								<template v-if="ban.warning">.</template>
								<template v-else>
									<template v-if="!ban.expires_at">
										<span class="bad">and is permanent.</span>
									</template>
									<template v-else>
										and expires at
										<b>{{ GetPostTimeReadable(ban.expires_at) }}</b>.
									</template>
								</template>

								<template v-if="expired(ban)">
									<br/>
									<span class="good">
										Now that you have seen this message, this
										<template v-if="ban.warning">warning</template><template v-else>ban</template>
										is no longer active.
									</span>
								</template>
							</p>
							<hr />
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
	.error {
		color: var(--user-error-color);
	}

	.bad {
		font-weight: bold;
		color: var(--user-error-color);
	}

	.good {
		font-weight: bold;
		color: var(--username-color);
	}


	.ban-info {
		.ban-reason {
			margin-left: 2em;
			white-space: pre;
		}
	}

	.container {
		border: 1px solid var(--banner-title-color);
		border-radius: 15px;
		width: 70%;
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
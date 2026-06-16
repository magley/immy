<script setup lang="ts">
	import { BanAPI, type BanDTO } from '@/api/ban.api';
	import type { ApiResponse } from '@/api/http';
	import type { AxiosError, AxiosResponse } from 'axios';
	import { onMounted, ref } from 'vue';

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
		}).catch((err: AxiosError) => {
			error.value = "Could not fetch your bans";
			console.log(err);
		}).finally(() => {
			loading.value = false;
		});
	}

</script>

<template>
	<div class="container">
		<div class="header">
			<h2>Banned?</h2>
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
						{{ bans }}
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
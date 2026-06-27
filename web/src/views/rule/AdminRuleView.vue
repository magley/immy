<script setup lang="ts">
	import { RuleAPI, type RuleDTO } from '@/api/rule.api';
	import { UserAPI, UserRole } from '@/api/user.api';
	import type { AxiosError } from 'axios';
	import { onMounted, reactive, ref } from 'vue';
	import { useRoute, useRouter } from 'vue-router';
	import PaginatorComponent from '@/components/PaginatorComponent.vue';
	import { Paginator } from '@/util/pagination.util';

	const route = useRoute();
	const router = useRouter();
	const error = ref<string | undefined>(undefined);
	const rules = ref<RuleDTO[]>([]);
	const paginator = reactive<Paginator<RuleDTO[]>>(new Paginator(RuleAPI.ListRules));

	onMounted(() => {
		UserAPI.AuthorizeUser({required_roles: [UserRole.Admin]}).then(() => {
			getRules();
		}).catch((err: AxiosError) => {
			console.error(err);
			router.push({ path: '/login', query: { redirect: route.fullPath } });
		});
	});

	const getRules = () => {
		paginator.getItems()
		.then((res) => rules.value = res.data.data!)
		.catch((_) => error.value = "Could not fetch Rules!");
	}

	const gotoPage = (p: number) => {
		paginator.page = p;
		getRules();
	}

	const removeRule = (Rule: RuleDTO) => {
		RuleAPI.DeleteRule(Rule.id).then((res) => {
			refreshRule(Rule.id);
		}).catch((err: AxiosError) => {
			error.value = `Failed to delete Rule #${Rule.id}`;
			console.error(err);
		});
	}

	const refreshRule = (RuleID: number) => {
		RuleAPI.GetRule(RuleID).then((res) => {
			for (let i = 0; i < rules.value.length; i++) {
				if (rules.value[i]!.id == RuleID) {
					rules.value[i] = res.data.data!
				}
			}
		}).catch((err: AxiosError) => {
			error.value = `Failed to refresh Rule #${RuleID}`;
			console.error(err);
		})
	}
</script>

<template>
	<h1>Rules</h1>

	<div class="error" v-if="error">{{ error }}</div>
	<div v-if="paginator.loading">Loading...</div>
	<div v-else>
		<div v-if="rules.length > 0">
			<div class="center nav">
				<PaginatorComponent :paginator="paginator" @goto-page="gotoPage" />
				<br/>
			</div>

			<table>
				<tbody>
					<tr>
						<th>ID</th>
						<th>Name</th>
						<th>Global?</th>
						<th>Description</th>
						<th>Danger</th>
					</tr>
					<template v-for="rule, i of rules">
						<tr>
							<td>{{ rule.id }}</td>
							<td class="center">{{ rule.is_global }}</td>
							<td>{{ rule.description }}</td>
							<td>{{ rule.danger }}</td>
						</tr>
					</template>
				</tbody>
			</table>

			<div class="center nav">
				<PaginatorComponent :paginator="paginator" @goto-page="gotoPage" />
			</div>
		</div>
		<div v-else class="center">
			No rules defined.
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
		color: var(--Rulener-title-color);
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
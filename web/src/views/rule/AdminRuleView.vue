<script setup lang="ts">
	import { RuleAPI, type RuleDTO, type CreateRuleDTO, type RuleBoardDTO } from '@/api/rule.api';
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
	const ruleBoards = ref<RuleBoardDTO[]>([]);
	const paginator = reactive<Paginator<RuleDTO[]>>(new Paginator(RuleAPI.ListRules));


	const newRule = ref<CreateRuleDTO>({
		title: '',
		description: '',
		is_global: false,
		danger: 0
	});
	const newRuleError = ref<string | undefined>(undefined);
	const newRuleFormVisible = ref<boolean>(false);

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
		.then((res) => {
			rules.value = res.data.data!;
			RuleAPI.ListAllRuleBoards().then((res) => {
				ruleBoards.value = res.data.data!
			}).catch((err) => {
				error.value = "Could not fetch Rules!";
				console.error(err);
			});
		}).catch((_) => error.value = "Could not fetch Rules!");
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

	const onClickCreateNewRule = () => {
		newRuleError.value = undefined;

		newRule.value.title = newRule.value.title.trim();
		newRule.value.description = newRule.value.description.trim();

		if (newRule.value.title == "") {
			newRuleError.value = "Title must not be empty";
			return;
		}
		if (newRule.value.description == "") {
			newRuleError.value = "Description must not be empty";
			return;
		}

		RuleAPI.CreateRule(newRule.value).then((res) => {
			getRules();

			newRule.value = {
				title: '',
				description: '',
				is_global: false,
				danger: 0
			};
		}).catch((err) => {
			newRuleError.value = "Could not create new rule";
			console.error(err);
		});
	}

	const toggleNewRuleFormVisible = () => {
		newRuleFormVisible.value = !newRuleFormVisible.value;
	}
</script>

<template>
	<h1>Rules</h1>

	<div class="error" v-if="error">{{ error }}</div>
	<div v-if="paginator.loading">Loading...</div>
	<div>
		<hr />
		<div>
			<div v-if="!newRuleFormVisible">
				<div class="center">[<a href='#' @click.prevent="toggleNewRuleFormVisible" >Create New Rule</a>]</div>
			</div>
			<div v-if="newRuleFormVisible">
				<div class="center">[<a href='#' @click.prevent="toggleNewRuleFormVisible">Hide Form</a>]</div>
				<div class="new-rule-form">
					<h3>Create new rule</h3>
					<label for="new-rule-title">Title:</label>
					<input id="new-rule-title" placeholder="SFW Content only" required v-model="newRule.title">
					<br/>
					<label for="new-rule-description">Description:</label><br/>
					<textarea id="new-rule-description" placeholder="All content and media uploaded by the user must be 'Safe For Work'." required v-model="newRule.description" cols="30" rows=5></textarea>
					<br/>
					<label for="new-rule-global"><abbr title="Global rules are applicable across ALL boards on the website">Is global</abbr>:</label>
					<input id="new-rule-global" type="checkbox" v-model="newRule.is_global">
					<br />
					<label for="new-rule-danger"><abbr title="Used internally as a priority value when posts are reported for violating rules. Greater values are more dangerous.">Danger level</abbr>:</label>
					<input id="new-rule-danger" type="number" v-model="newRule.danger">
					<br/>
					<button @click="onClickCreateNewRule" class="submit">Submit</button>
					<div v-if="newRuleError" class="error">{{ newRuleError }}</div>
				</div>
			</div>
			<hr />
		</div>

		<div v-if="rules.length > 0">
			<div class="center nav">
				<PaginatorComponent :paginator="paginator" @goto-page="gotoPage" empty-message="There are no rules defined yet" />
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
							<td>{{ rule.title }}</td>
							<td class="center">{{ rule.is_global }}</td>
							<td>{{ rule.description }}</td>
							<td>{{ rule.danger }}</td>
						</tr>
					</template>
				</tbody>
			</table>

			<div class="center nav">
				<PaginatorComponent :paginator="paginator" @goto-page="gotoPage" empty-message="There are no rules defined yet" />
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

	h1, h2, h3 {
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

	.new-rule-form {
		width: 40%;
		padding: 2em 1em;
		border: 1px solid black;
		margin: auto;

		input:not([type]), input[type="text"], textarea {
			width: 100%;
		}

		button.submit {
			text-align: right;
		}
	}
</style>
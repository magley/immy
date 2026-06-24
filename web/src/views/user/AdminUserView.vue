<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { UserAPI, type UserDTO, type CreateUserDTO, type UpdateUserDTO, UserRole} from "@/api/user.api.ts";
	import type { ApiResponse } from '@/api/http';
	import type { AxiosResponse, AxiosError } from 'axios';
	import { useRoute, useRouter } from 'vue-router';

	const router = useRouter();
	const route = useRoute();

	const users = ref<UserDTO[]>([]);
	const createUserDTO = ref<CreateUserDTO>({
		role: UserRole.Janitor,
		username: '',
		password: ''
	});
	const createUserError = ref<string>("");
	
	onMounted(() => {
		UserAPI.AuthorizeUser({required_roles: [UserRole.Admin]}).then(() => {
			get_users();
		}).catch((err: AxiosError) => {
			console.error(err);
			router.push({ path: '/login', query: { redirect: route.fullPath } });
		});
	});
	
	const get_users = () => {
		UserAPI.ListUsers().then((res: AxiosResponse<ApiResponse<UserDTO[]>>) => {
			users.value = res.data.data!;
		}).catch((err) => {
			console.error(err);
		});
	}
	
	const onSubmitCreateUser = () => {
		createUserError.value = "";
		
		UserAPI.CreateUser(createUserDTO.value).then((res: AxiosResponse<ApiResponse<UserDTO>>) => {
			get_users();
		}).catch((err: AxiosError) => {
			createUserError.value = "Could not create user";
			console.error(err);
		});
	}
	
	const onSubmitChangesToUser = (idx: number) => {
		createUserError.value = "";
		
		const user: UserDTO = users.value[idx]!;
		const updateDto: UpdateUserDTO = {
			role: user.role,
			username: null
		}
		
		UserAPI.UpdateUser(user.id, updateDto).then((res: AxiosResponse<ApiResponse<UserDTO>>) => {
			get_users();
		}).catch((err: AxiosError) => {
			createUserError.value = "Could not update user";
			console.error(err);
		});
	}
	
	const onDeleteUser = (idx: number) => {
		const user: UserDTO = users.value[idx]!;
		
		UserAPI.DeleteUser(user.id).then((res: AxiosResponse<ApiResponse<number>>) => {
			get_users();
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}
</script>

<template>
	<h1>Users</h1>
	
	<form @submit.prevent="onSubmitCreateUser">
		<input type=text placeholder="Username" required v-model="createUserDTO.username"/>
		<br/>
		<input type=password placeholder="Password" required v-model="createUserDTO.password"/>
		<br/>
		<select v-model="createUserDTO.role">
			<option disabled value="">Please select one</option>
			<option>{{UserRole.Admin}}</option>
			<option>{{UserRole.Moderator}}</option>
			<option>{{UserRole.Janitor}}</option>
		</select>
		<br/>
		<button type=submit>Create user</button>
		
		<template v-if="createUserError">
			<div/>
			<span class="error">{{createUserError}}</span>
		</template>
	</form>
	
	<br />

	<table>
		<tbody>
			<tr>
				<th>ID</th>
				<th>Username</th>
				<th>Role</th>
				<th>Update</th>
				<th>Delete</th>
			</tr>
			<tr v-for="user, i in users">
				<td>{{user.id}}</td>
				<td><img :src="`/icons/user-role-${user.role}.gif`" :title="user.role" /> {{user.username}}</td>
				<td>
					<select v-model="user.role">
						<option disabled value="">Please select one</option>
						<option>{{UserRole.Admin}}</option>
						<option>{{UserRole.Moderator}}</option>
						<option>{{UserRole.Janitor}}</option>
					</select>
				</td>
				<td>
					<button @click="onSubmitChangesToUser(i)">Update</button>
				</td>
				<td>
					<button @click="onDeleteUser(i)">Delete</button>
				</td>
			</tr>
		</tbody>
	</table>
</template>

<style scoped>
	h1 {
		color: var(--banner-title-color);
	}

	table {
		border: 1px solid black;
	}
	
	th {
		background-color: var(--background-color-darker);
		padding: 0em 1em;
	}
	
	.error {
		color: var(--user-error-color);
	}
</style>
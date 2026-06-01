<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { UserAPI, UserDTO, CreateUserDTO, UpdateUserDTO, UserType, LoginUserDTO, JWTClaims} from "@/api/user.api.ts";

	const users = ref<UserDTO[]>([]);
	const createUserDTO = ref<CreateUserDTO>({type:UserType.Janitor});
	const createUserError = ref<string>(undefined);
	
	onMounted(() => {
		get_users();
	});
	
	const get_users = () => {
		UserAPI.ListUsers().then((res: AxiosResponse<ApiResponse<UserDTO[]>>) => {
			users.value = res.data.data;
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
	
	const onSubmitChangesToUser = (idx: int) => {
		createUserError.value = "";
		
		const user: UserDTO = users.value[idx];
		const updateDto: UpdateUserDTO = {
			type: user.type
		}
		
		UserAPI.UpdateUser(user.id, updateDto).then((res: AxiosResponse<ApiResponse<UserDTO>>) => {
			get_users();
		}).catch((err: AxiosError) => {
			createUserError.value = "Could not update user";
			console.error(err);
		});
	}
	
	const onDeleteUser = (idx: int) => {
		const user: UserDTO = users.value[idx];
		
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
		<select v-model="createUserDTO.type">
			<option disabled value="">Please select one</option>
			<option>{{UserType.Admin}}</option>
			<option>{{UserType.Moderator}}</option>
			<option>{{UserType.Janitor}}</option>
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
				<td>{{user.username}}</td>
				<td>
					<select v-model="user.type">
						<option disabled value="">Please select one</option>
						<option>{{UserType.Admin}}</option>
						<option>{{UserType.Moderator}}</option>
						<option>{{UserType.Janitor}}</option>
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
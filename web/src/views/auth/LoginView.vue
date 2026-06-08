<script setup lang="ts">
	import { ref } from 'vue';
	import { UserAPI, type LoginUserDTO, type LoginResponseDTO} from "@/api/user.api.ts";
	import type { ApiResponse } from '@/api/http';
	import type { AxiosResponse, AxiosError } from 'axios';
import BoardListNav from '@/components/board/BoardListNav.vue';

	const loginUserDTO = ref<LoginUserDTO>({
		username: '',
		password: ''
	});
	const errorMessage = ref<string | undefined>(undefined);

	const onSubmitLogin = () => {
		errorMessage.value = "";
		UserAPI.LoginUser(loginUserDTO.value).then((res: AxiosResponse<ApiResponse<LoginResponseDTO>>) => {
			const dto: LoginResponseDTO = res.data.data!;
			onSuccess(dto);
		}).catch((err: AxiosError) => {
			if (err.status! == 401) {
				errorMessage.value = "Incorrect username or password";
			} else {
				errorMessage.value = "Failed to log in";
			}
			console.error(err);
		});
	}
	
	const onSuccess = (dto: LoginResponseDTO) => {
		localStorage.setItem("username", dto.username);
		localStorage.setItem("id", `${dto.id}`);
		localStorage.setItem("role", dto.type);
		localStorage.setItem("jwt", dto.jwt);
		location.reload();
	}
</script>

<template>
	<BoardListNav :isCatalog=false />
	<hr/>

	<form @submit.prevent="onSubmitLogin">
		<h1>Login</h1>
		
		<input type=text placeholder="Username" required v-model="loginUserDTO.username"/>
		<br/>
		<input type=password placeholder="Password" required v-model="loginUserDTO.password"/>
		<br/>
		<br/>
		<button type=submit>Sign In</button>
		
		<template v-if="errorMessage">
			<div/>
			<span class="error">{{errorMessage}}</span>
		</template>
	</form>

	<hr/>
	<BoardListNav :isCatalog=false />
</template>

<style scoped>
	form {
		width: 40%;
		margin: 2em auto;
		border: 1px solid black;
		background-color: var(--background-color-darker);
		padding: 1em;
		text-align: center;
	}

	h1 {
		color: var(--banner-title-color);
	}
	.error {
		color: var(--user-error-color);
	}
</style>
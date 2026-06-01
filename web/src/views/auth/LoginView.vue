<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { UserAPI, UserDTO, CreateUserDTO, UpdateUserDTO, UserType, LoginUserDTO, JWTClaims} from "@/api/user.api.ts";

	const loginUserDTO = ref<LoginUserDTO>({});
	const errorMessage = ref<string>(undefined);

	const onSubmitLogin = () => {
		errorMessage.value = "";
		UserAPI.LoginUser(loginUserDTO.value).then((res: AxiosResponse<ApiResponse<LoginResponseDTO>>) => {
			const dto: LoginResponseDTO = res.data.data;
			onSuccess(dto);
		}).catch((err: AxiosError) => {
			onFail(err);
		});
	}
	
	const onSuccess = (dto: LoginResponseDTO) => {
		localStorage.setItem("username", dto.username);
		localStorage.setItem("id", dto.id);
		localStorage.setItem("role", dto.type);
		localStorage.setItem("jwt", dto.jwt);
	}
	
	const onFail = (error: AxiosError) => {
		if (error.status == 401) {
			errorMessage.value = "Incorrect username or password";
		} else if (error.ok) {
			errorMessage.value = "Failed to log in";
			console.error(error.response.data);
		}
	}
</script>

<template>
	<form @submit.prevent="onSubmitLogin">
		<h2>Login</h2>
		
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
</template>

<style scoped>
	.error {
		color: var(--user-error-color);
	}
</style>
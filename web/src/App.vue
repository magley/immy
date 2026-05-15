<script setup lang="ts">
	import { BoardDTO, BoardAPI } from '@/api/board.api.ts';
	import { ref, onMounted } from 'vue';
	
	const boards = ref<BoardDTO[]>([]);
	const boardsError = ref<string | null>(null);
	
	onMounted(() => {
		BoardAPI.ListBoards().then((res: AxiosResponse<ApiResponse<BoardDTO[]>>) => {
			boards.value = res.data.data;
		}).catch((err: AxiosError) => {
			boardsError.value = "Failed to fetch boards";
			console.error(err);	
		});
	});
	
</script>


<template>
	<nav>
		<RouterLink to="/"><b>ImmyChan</b></RouterLink> 
		|
		<RouterLink to="/login">Log In</RouterLink>
	</nav>
	
	<nav>
		<b>[Admin]</b>
		|
		<RouterLink to="/admin-users">Users</RouterLink>
		|
		<RouterLink to="/admin-boards">Boards</RouterLink>
	</nav>
	
	<nav>
		<template v-if="boardsError">
			{{ boardsError }}
		</template>
		<template v-else>
			[
			<template v-for="board, i in boards">
				<RouterLink :to="`/${board.code}`">{{ board.code}}</RouterLink>
				<template v-if="i != boards.length - 1"> / </template>
			</template>
			]
		</template>
	</nav>
	
	<p><strong>Current route path:</strong> {{ $route.fullPath }}</p>
	
	<main>	
		<!-- The key is so the page resets when the route changes. Don't use
		<!-- fullPath because then the '#abc' anchor will cause a reset as
		<!-- well. -->
		<RouterView :key="$route.path"/>
	</main>
</template>


<!-- 
<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { BoardAPI, BoardDTO, CreateBoardDTO, UpdateBoardDTO } from "./api/board.api.ts";
	import { UserAPI, UserDTO, CreateUserDTO, UpdateUserDTO, UserType, LoginUserDTO, JWTClaims} from "./api/user.api.ts";

	const boards = ref<BoardDTO[]>([]);
	const users = ref<UserDTO[]>([]);
	const createUserDTO = ref<CreateUserDTO>({type:UserType.Janitor});
	const loginUserDTO = ref<LoginUserDTO>({});

	const get_boards = () => {
		BoardAPI.ListBoards().then((res: AxiosResponse<ApiResponse<BoardDTO[]>>) => {
			boards.value = res.data.data;
		}).catch((err) => {
			console.error(err);
		});
	};
	

	const get_users = () => {
		UserAPI.ListUsers().then((res: AxiosResponse<ApiResponse<UserDTO[]>>) => {
			users.value = res.data.data;
		}).catch((err) => {
			console.error(err);
		});
	}
	
	const onSubmitCreateUser = () => {
		UserAPI.CreateUser(createUserDTO.value).then((res: AxiosResponse<ApiResponse<UserDTO>>) => {
			get_users();
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}
	
	const onSubmitLoginUser = () => {
		UserAPI.LoginUser(loginUserDTO.value).then((res: AxiosResponse<ApiResponse<LoginResponseDTO>>) => {
			const dto: LoginResponseDTO = res.data.data;
			localStorage.setItem("username", dto.username);
			localStorage.setItem("id", dto.id);
			localStorage.setItem("role", dto.type);
			localStorage.setItem("jwt", dto.jwt);
		}).catch((err: AxiosError) => {
			console.error(err);
		});
	}
	
	onMounted(() => {
		get_boards();
		get_users();
	});
	
</script>

<template>
	<h1>Boards</h1>
	<ul>
		[
		<span v-for="board in boards">
			<a :href='`/${board.code}`'>/{{ board.code }}/</a> ♦
		</span>
		]
	</ul>
	<h1>Users</h1>
	<ul>
		<li v-for="user in users">
			{{user.username}} - {{ user.type }}
		</li>
	</ul>
	
	<form @submit.prevent="onSubmitCreateUser">
		<h2>Create</h2>
		
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
	</form>
	
	<form @submit.prevent="onSubmitLoginUser">
		<h2>Login</h2>
		
		<input type=text placeholder="Username" required v-model="loginUserDTO.username"/>
		<br/>
		<input type=password placeholder="Password" required v-model="loginUserDTO.password"/>
		<br/>
		<br/>
		<button type=submit>Sign In</button>
	</form>
</template>

<style scoped></style>
 -->
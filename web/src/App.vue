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
	<span id="top"></span>
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
		fullPath because then the '#abc' anchor will cause a reset as
		ell. -->
		<RouterView :key="$route.path"/>
	</main>
	
	<span id="bottom"></span>
</template>
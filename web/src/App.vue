<script setup lang="ts">
	import { ref, onMounted } from 'vue';
	import { UserRole } from './api/user.api';

	const userRole = ref<string | undefined>(undefined);
	const userName = ref<string | undefined>(undefined);
	const theme = ref<string>("yotsuba");

	onMounted(() => {
		userRole.value = localStorage.getItem("role") ?? undefined;
		userName.value = localStorage.getItem("username") ?? undefined;

		setTheme(localStorage.getItem("theme") ?? "yotsuba");
		theme.value = localStorage.getItem("theme") ?? "yotsuba";
	});

	const themeChanged = () => {
		setTheme(theme.value);
	}

	const setTheme = (theme: string) => {
		document.documentElement.setAttribute('data-theme', theme);
		localStorage.setItem("theme", theme);

		const codeThemeDict: Record<string, string> = {
			'yotsuba': 'a11y-light',
			'futaba': 'a11y-light',

			'yotsuba-b': 'stackoverflow-light',
			'burichan': 'stackoverflow-light',
		}
		const codeTheme: string = codeThemeDict[theme] ?? 'a11y-light';

		const hljsThemeLink: HTMLLinkElement = document.getElementById('code-syntax-theme')! as any as HTMLLinkElement;
		hljsThemeLink.href = `https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/${codeTheme}.min.css`;
	}

	const logOut = () => {
		localStorage.removeItem("username");
		localStorage.removeItem("id");
		localStorage.removeItem("role");
		localStorage.removeItem("jwt");

		location.reload();
	}
</script>


<template>
	<span id="top"></span>
	<nav>
		<RouterLink to="/"><b>ImmyChan</b></RouterLink> 
		|
		<RouterLink to="/login">Log In</RouterLink>
	</nav>

	
	<nav v-if="userRole != undefined">
		<b><img :src="`/icons/user-role-${userRole}.gif`" :title="userRole" class="icon" />
			{{userName}}
			<span class="capcode" :class="userRole" >## {{ userRole }}</span></b>
		|
		<a href="#" @click.prevent="logOut">Log Out</a>
		|
		<span v-if="userRole == UserRole.Admin">
			<RouterLink to="/admin-users">Users</RouterLink>
			|
			<RouterLink to="/admin-boards">Boards</RouterLink>
			|
			<RouterLink to="/admin-bans">Bans</RouterLink>
			|
		</span>
		<span v-if="userRole == UserRole.Moderator">
			<RouterLink to="/admin-bans">Bans</RouterLink>
			|
		</span>
		<span v-if="userRole == UserRole.Janitor">

		</span>
	</nav>

	<main>	
		<!-- The key is so the page resets when the route changes. Don't use
		fullPath because then the '#abc' anchor will cause a reset as
		ell. -->
		<RouterView :key="$route.path"/>
	</main>

	<label>Theme:</label>
	<select v-model="theme" @change="themeChanged">
		<option value="yotsuba">Yotsuba</option>
		<option value="yotsuba-b">Yotsuba B</option>
		<option value="futaba">Futaba</option>
		<option value="burichan">Burichan</option>
	</select>
	
	<span id="bottom"></span>
</template>

<style scoped>
</style>
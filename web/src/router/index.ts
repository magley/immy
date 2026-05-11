import { createRouter, createWebHistory } from 'vue-router';
import HomeView from "@/views/HomeView.vue";
import LoginView from "@/views/auth/LoginView.vue";
import UserView from "@/views/user/UserView.vue";


const routes = [
	{ path: '/', component: HomeView },
	{ path: '/login', component: LoginView },
	{ path: '/users', component: UserView },
]

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: routes,
})

export default router

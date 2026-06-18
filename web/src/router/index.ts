import { createRouter, createWebHistory } from 'vue-router';
import HomeView from "@/views/HomeView.vue";
import LoginView from "@/views/auth/LoginView.vue";
import AdminUserView from "@/views/user/AdminUserView.vue";
import AdminBoardView from "@/views/board/AdminBoardView.vue";
import BoardHomeView from "@/views/board/BoardHomeView.vue";
import BoardCatalogView from "@/views/board/BoardCatalogView.vue";
import ThreadView from "@/views/thread/ThreadView.vue";
import BoardArchiveView from '@/views/board/BoardArchiveView.vue';
import BannedView from '@/views/ban/BannedView.vue';
import AdminBanView from '@/views/ban/AdminBanView.vue';

const routes = [
	{ path: '/', component: HomeView },
	{ path: '/login', component: LoginView },
	{ path: '/banned', component: BannedView },
	{ path: '/admin-users', component: AdminUserView },
	{ path: '/admin-boards', component: AdminBoardView },
	{ path: '/admin-bans', component: AdminBanView },
	{ path: '/:board_code/catalog', component: BoardCatalogView },
	{ path: '/:board_code/archive', component: BoardArchiveView },
	{ path: '/:board_code/thread/:thread_num', component: ThreadView },
	{ path: '/:board_code', component: BoardHomeView },
]

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: routes,
})

export default router

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
import PublicBansView from '@/views/ban/PublicBansView.vue';
import BlogpostView from '@/views/blogpost/BlogpostView.vue';
import SingleBlogpostView from '@/views/blogpost/SingleBlogpostView.vue';
import AdminRuleView from '@/views/rule/AdminRuleView.vue';

const routes = [
	{ path: '/', component: HomeView },
	{ path: '/login', component: LoginView },
	{ path: '/banned', component: BannedView },
	{ path: '/bans', component: PublicBansView },
	{ path: '/blog', component: BlogpostView },
	{ path: '/blog/:id', component: SingleBlogpostView },
	{ path: '/admin-users', component: AdminUserView },
	{ path: '/admin-boards', component: AdminBoardView },
	{ path: '/admin-bans', component: AdminBanView },
	{ path: '/admin-rules', component: AdminRuleView },
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

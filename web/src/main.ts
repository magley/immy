import { createApp } from 'vue';
import App from '@/App.vue';
import router from './router';
import "../public/main.css";
import VueLatex from 'vatex';

const app = createApp(App);
app.use(router);
app.use(VueLatex);
app.mount('#app');
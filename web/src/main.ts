import { createApp } from 'vue';
import App from '@/App.vue';
import router from './router';
import "../public/main.css";
import VueLatex from 'vatex';

import hljs from 'highlight.js/lib/core';
import javascript from 'highlight.js/lib/languages/javascript';
import hljsVuePlugin from "@highlightjs/vue-plugin";

hljs.registerLanguage('javascript', javascript);

const app = createApp(App);
app.use(router);
app.use(VueLatex);
app.use(hljsVuePlugin);
app.mount('#app');
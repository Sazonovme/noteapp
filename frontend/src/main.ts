import { createPinia } from 'pinia';
import { createApp } from 'vue';

import 'element-plus/dist/index.css';
import 'element-plus/theme-chalk/dark/css-vars.css';

import { App } from '@app/index';

import { router } from '@pages/routing';

import '@shared/styles/index.css';

createApp(App)
    .use(createPinia())
    .use(router)
    .mount('#app');

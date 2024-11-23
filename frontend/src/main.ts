import { createApp } from 'vue';

import { App } from 'app/index';

import { router } from 'pages/routing';

import 'shared/styles/index.css';

createApp(App)
    .use(router)
    .mount('#app');

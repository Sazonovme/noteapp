import { createMemoryHistory, createRouter } from 'vue-router';

import { PageHome } from './page-home';

const routes = [
    { path: '/', component: PageHome },
];

export const router = createRouter({
    history: createMemoryHistory(),
    routes,
});

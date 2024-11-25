import { createMemoryHistory, createRouter } from 'vue-router';

import { PageHome } from './page-home';
import { PageAuthorization } from './page-authorization';

const routes = [
    {
        path: '/',
        name: 'home',
        component: PageHome,
        meta: {
            layout: 'base',
        },
    },
    {
        path: '/auth',
        // alias: '/auth',
        name: 'auth',
        component: PageAuthorization,
        meta: {
            layout: 'base',
        },
    },
];

export const router = createRouter({
    history: createMemoryHistory(),
    routes,
});

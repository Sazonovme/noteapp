import { createMemoryHistory, createRouter } from 'vue-router';

import { ROUTES_PATH_COMMON } from '@shared/constants';

import { PageHome } from './page-home';
import { PageAuthorization } from './page-authorization';

const routes = [
    {
        path: ROUTES_PATH_COMMON.HOME,
        name: 'home',
        component: PageHome,
        // component: () => import('./page-home/page-home.vue'),
        meta: {
            layout: 'base',
        },
    },
    {
        path: ROUTES_PATH_COMMON.AUTH,
        // alias: '/auth',
        name: 'auth',
        component: PageAuthorization,
        // component: () => import('./page-authorization/page-authorization.vue'),
        meta: {
            layout: 'auth',
        },
    },
];

export const router = createRouter({
    history: createMemoryHistory(),
    routes,
});

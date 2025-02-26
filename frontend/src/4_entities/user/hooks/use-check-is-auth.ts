import { useRouter } from 'vue-router';
import { onMounted, ref } from 'vue';

import { LOCAL_STORAGE_UPDATED_EVENT, type EventStorageUpdatedTypes } from '@shared/utils';
import { ROUTES_PATH_COMMON } from '@shared/constants';

import { getAuthTokens } from '../utils';
import { ACCESS_TOKEN } from '../constants';

export const useCheckIsAuth = () => {
    const { accessToken } = getAuthTokens();
    const isAuth = ref(!!accessToken);
    const nav = useRouter();

    const changeTokens = (event: Event) => {
        const { key, value } = event as unknown as EventStorageUpdatedTypes;
        if (key === ACCESS_TOKEN) {
            isAuth.value = !!value;
            if (!value) {
                nav.push(ROUTES_PATH_COMMON.AUTH);
            }
        }
    };

    onMounted(() => {
        window.addEventListener(LOCAL_STORAGE_UPDATED_EVENT, changeTokens, false);

        return () => window.removeEventListener(LOCAL_STORAGE_UPDATED_EVENT, changeTokens, false);
    });

    return {
        isAuth: isAuth.value,
    };
};

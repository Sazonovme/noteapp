import { storage } from '@shared/utils';

import { ACCESS_TOKEN, REFRESH_TOKEN, TYPE_TOKEN } from '../constants';

export const setAuthTokens = ({ access, refresh, type }: { access: string, refresh: string, type: string; }) => {
    storage.local.setItem(ACCESS_TOKEN, access);
    storage.local.setItem(REFRESH_TOKEN, refresh);
    storage.local.setItem(TYPE_TOKEN, type || 'Bearer');
};

export const getAuthTokens = () => ({
    accessToken: storage.local.getItem(ACCESS_TOKEN),
    refreshToken: storage.local.getItem(REFRESH_TOKEN),
    typeToken: storage.local.getItem(TYPE_TOKEN),
});

export const clearAuthTokens = () => {
    storage.local.removeItem(ACCESS_TOKEN);
    storage.local.removeItem(REFRESH_TOKEN);
    storage.local.removeItem(TYPE_TOKEN);
};

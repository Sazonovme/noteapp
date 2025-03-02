import { storage } from '@shared/utils';

import { ACCESS_TOKEN, REFRESH_TOKEN } from '../constants';

export const setAuthTokens = ({ access, refresh }: { access: string, refresh: string }) => {
    storage.local.setItem(ACCESS_TOKEN, access);
    storage.local.setItem(REFRESH_TOKEN, refresh);
};

export const getAuthTokens = () => ({
    accessToken: storage.local.getItem(ACCESS_TOKEN),
    refreshToken: storage.local.getItem(REFRESH_TOKEN),
});

export const clearAuthTokens = () => {
    storage.local.removeItem(ACCESS_TOKEN);
    storage.local.removeItem(REFRESH_TOKEN);
};

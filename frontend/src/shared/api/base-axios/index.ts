import axios, { type AxiosError, type AxiosRequestConfig, type AxiosPromise, type InternalAxiosRequestConfig } from 'axios';

import { clearAuthTokens, getAuthTokens, setAuthTokens } from '@entities/user';

import { API_BASE_URL } from '@shared/constants';

import api from '../index';

declare module 'axios' {
    export interface AxiosRequestConfig {
        includeAuthToken?: boolean;
    }
}

const instance = axios.create({
    baseURL: API_BASE_URL,
});

const _setHeaders = (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    if (config.includeAuthToken !== false) {
        const { typeToken, accessToken } = getAuthTokens();
        config.headers.Authorization = `${typeToken} ${accessToken}`;
    }

    return config;
};

const _refreshToken = async (config: InternalAxiosRequestConfig) => {
    try {
        const { refreshToken } = getAuthTokens();

        if (!refreshToken) {
            clearAuthTokens();
            // eslint-disable-next-line prefer-promise-reject-errors
            return Promise.reject({ data: { message: 'refresh_token отсутствует' }, status: 401 });
        }

        const { data } = await api.authorization.refreshToken({ refresh_token: refreshToken });
        setAuthTokens({ access: data.access_token, refresh: data.refresh_token, type: data.token_type });

        return instance(config);
    } catch (error) {
        console.log('Ошибка при получении refresh_token: ', error);
        clearAuthTokens();
        // eslint-disable-next-line prefer-promise-reject-errors
        return Promise.reject({ data: { message: 'Ошибка при получении refresh_token: ' }, status: 401 });
    }
};

instance.interceptors.request.use(
    config => _setHeaders(config),
    async error => Promise.reject(error)
);

instance.interceptors.response.use(
    response => response,
    async error => {
        if (error.response?.status === 401) return _refreshToken(error.config);
        return Promise.reject(error.response);
    }
);

export { instance as axios, type AxiosPromise, type AxiosError, type AxiosRequestConfig };

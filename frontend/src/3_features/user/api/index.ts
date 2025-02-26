import { axios, AxiosPromise } from '@shared/api';

export const login = (val: { email: string, password: string }): AxiosPromise<{ accessToken: string, refreshToken: string }> => axios.post('/sign-in', { data: val });

export const logout = () => axios.get('/logout');

export const registration = (val: { email: string, password: string }): AxiosPromise<{ accessToken: string, refreshToken: string }> => axios.post('/sign-up', { data: val });

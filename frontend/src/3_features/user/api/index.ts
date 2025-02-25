import { axios } from '@shared/api';

export const login = (val: { email: string, password: string }): Promise<{ accessToken: string, refreshToken: string }> => axios.post('/sign-in', val);

export const logout = () => axios.get('/logout');

export const registration = (val: { email: string, password: string }): Promise<{ accessToken: string, refreshToken: string }> => axios.post('/sign-up', val);

import { axios } from '@shared/api';

export const login = (val: Record<string, string>) => axios.post('/management/auth/login', val);

export const logout = () => axios.get('/management/oauth2/logout');

export const registration = (val: Record<string, string>) => axios.post('/management/oauth2/logout', val);

import { axios, type AxiosPromise } from '../../base-axios';

export type RefreshReqDto = {
    refreshToken: string,
};

export type RefreshResDto = {
    accessToken: string,
    refreshToken: string,
};

export const refreshToken = async (params: RefreshReqDto): AxiosPromise<RefreshResDto> => axios.post('/refresh-token', params, { includeAuthToken: false });

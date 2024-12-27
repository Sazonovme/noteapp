import { axios, type AxiosPromise } from '../../base-axios';

export type RefreshReqDto = {
    client_id?: 'management_panel',
    grant_type?: 'refresh_token',
    refresh_token: string,
};

export type RefreshResDto = {
    access_token: string,
    refresh_token: string,
    id_token: string,
    expires_in: number,
    token_type: string,
};

export const refreshToken = async (params: RefreshReqDto): AxiosPromise<RefreshResDto> => {
    const queryParams = new URLSearchParams({
        client_id: 'management_panel',
        grant_type: 'refresh_token',
        ...params,
    }).toString();

    return axios.post(`/oauth2/token?${queryParams}`, params, { includeAuthToken: false });
};

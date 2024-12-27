// import { baseApi, TagTypes } from 'shared/api';
// import { PUBLIC_API_V1 } from 'shared/constants';

// import type { GetUserInfoResDto } from '../model';

// export const userEntitiesApi = baseApi.injectEndpoints({
//     endpoints: build => ({
//         getUserInfo: build.query<GetUserInfoResDto, void>({
//             providesTags: [TagTypes.userInfo],
//             query: () => `${PUBLIC_API_V1}/userinfo`,
//         }),
//     }),
//     overrideExisting: false,
// });

// export const {
//     useGetUserInfoQuery,
// } = userEntitiesApi;

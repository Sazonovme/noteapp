import { axios, AxiosPromise } from '@shared/api';

import { preparedTreeList } from '../utils/prepared-tree-list';
import { NodeInfoType, PreparedTreeListType, ResponseTreeListType } from '../model';

export const getTreeList = async (): Promise<PreparedTreeListType> => {
    const resp = (await axios.get('/getNotesList')) as ResponseTreeListType;

    return preparedTreeList(resp);
};

export const getNote = (id: string): AxiosPromise<NodeInfoType> => {
    const query = new URLSearchParams({ id }).toString();

    return axios.get(`/getNote?${query}`);
};

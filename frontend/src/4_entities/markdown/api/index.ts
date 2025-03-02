import { axios } from '@shared/api';

import { preparedTreeList } from '../utils/prepared-tree-list';
import { PreparedTreeListType, ResponseTreeListType } from '../model';

export const getTreeList = async (): Promise<PreparedTreeListType> => {
    // const json = {
    //     groups: [
    //         {
    //             group_id: 'folder1',
    //             group_name: 'folder1',
    //             groups: [
    //                 {
    //                     group_id: 'folder1',
    //                     group_name: 'folder1',
    //                     groups: [
    //                         {
    //                             group_id: 'folder1',
    //                             group_name: 'folder1',
    //                             groups: [],
    //                             notes: [
    //                                 { note_id: 'note1', note_title: 'note1', note_text: 'notetext1' },
    //                             ],
    //                         },
    //                     ],
    //                     notes: [
    //                         { note_id: 'note1', note_title: 'note1', note_text: 'notetext1' },
    //                     ],
    //                 },
    //             ],
    //             notes: [
    //                 { note_id: 'note1', note_title: 'note1', note_text: 'notetext1' },
    //             ],
    //         },
    //     ],

    //     notes: [
    //         { note_id: 'note1', note_title: 'note1', note_text: 'notetext1' },
    //     ],
    // };

    const resp = (await axios.get('/getNotesList')) as ResponseTreeListType;

    return preparedTreeList(resp);
    // return new Promise((resolve, reject) => setTimeout(() => resolve(preparedTreeList({ data: json })), 2000));
};

export const getNote = (id: string): Promise<string> => {
    const query = new URLSearchParams({ id }).toString();

    return axios.get(`/getNote?${query}`);
    // return new Promise((resolve, reject) => setTimeout(() => resolve('new TEXT from AXIOS'), 2000));
};

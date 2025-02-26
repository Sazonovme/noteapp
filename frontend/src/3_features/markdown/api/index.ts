import { axios } from '@shared/api';

export const createNote = (val: { group_id: string, /* | '0', */ title: string }) => axios.post('/addNote', val);

// по кнопке + выделить title
export const updateNote = (val: { id: string, group_id?: string, text?: string, title?: string }) => axios.put('/updateNote', val);

export const deleteNote = (id: string) => {
    const query = new URLSearchParams({ id }).toString();

    return axios.delete(`/delNote?${query}`);
};

// dropdown
export const updateGroup = (val: { id: string, name: string } | { id: string, parentGroupId: string }) => axios.put('/updateGroup', val);

export const deleteGroup = (id: string) => {
    const query = new URLSearchParams({ id }).toString();

    return axios.delete(`/delGroup?${query}`);
};

export const createGroup = (val: { name: string, parentIdGroup?: string }) => axios.post('/addGroup', val);

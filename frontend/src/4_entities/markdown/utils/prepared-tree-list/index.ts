import { ResponseTreeListType, PreparedTreeListType } from '../../model';

const preparedNotes = (tree: ResponseTreeListType['data']['notes']): PreparedTreeListType => tree.map(val => ({
    id: val.note_id,
    title: val.note_title,
    isFolder: false,
    isNote: true,
}));

const preparedGroup = (tree: ResponseTreeListType['data']['groups']): PreparedTreeListType => tree.map(val => ({
    id: val.group_id,
    title: val.group_name,
    isFolder: true,
    isNote: false,
    children: [...preparedGroup(val.groups), ...preparedNotes(val.notes)] as [],
}));

export const preparedTreeList = (treeList: ResponseTreeListType): PreparedTreeListType => ([
    {
        id: '-1',
        title: 'Root',
        isFolder: true,
        isNote: false,
        children: [...preparedGroup(treeList.data.groups), ...preparedNotes(treeList.data.notes)] as [],
    },
]);

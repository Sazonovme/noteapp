export type ResponseTreeListNoteType = {
    note_id: string,
    note_title: string,
    // note_text: string,
};

export type ResponseTreeListGroupType = {
    group_id: string,
    group_name: string,
    groups: ResponseTreeListGroupType[],
    notes: ResponseTreeListNoteType[],
};

export type ResponseTreeListType = {
    data: {
        groups: ResponseTreeListGroupType[],
        notes: ResponseTreeListNoteType[],
    },
};

export type PreparedTreeListType = {
    id: string,
    title: string,
    children?: PreparedTreeListType[]
    isFolder: boolean,
    isNote: boolean,
}[];

import { ref, onBeforeMount } from 'vue';

import { createGroup, createNote, deleteGroup, updateGroup } from '@features/markdown/api/';

import { useMarkdownStore } from '@entities/markdown';
import { getNote, getTreeList } from '@entities/markdown/api';

const useState = () => {
    const isTreeDialogOpen = ref(false);

    return {
        isTreeDialogOpen,
    };
};

const useActions = (state: ReturnType<typeof useState>, store: ReturnType<typeof useMarkdownStore>) => {
    const onOpenTreeDialog = () => {
        state.isTreeDialogOpen.value = true;
    };

    const onCloseTreeDialog = () => {
        state.isTreeDialogOpen.value = false;
    };

    const onCreateNewFolder = async (id: string, e: any) => {
        e.preventDefault();
        e.stopPropagation();

        await createGroup({ name: 'New Folder', parentIdGroup: String(id) });
        onCloseTreeDialog();
    };

    const onCreateNewNote = async (id: string, e: any) => {
        e.preventDefault();
        e.stopPropagation();

        await createNote({ group_id: id, title: 'New Note' });
        onCloseTreeDialog();
    };

    const onDeleteFolder = async (id: string, e: any) => {
        e.preventDefault();
        e.stopPropagation();

        await deleteGroup(id);
        onCloseTreeDialog();
    };

    const onDropFile = async (draggingNode: any, dropNode: any) => {
        //todo возможно разные запросы для заметок и для папок
        await updateGroup({ id: draggingNode.id, newIdGroup: dropNode.id });
    };

    const onOpenNote = async (id: string) => {
        store.setCurrentOpenIdNode(id);
        const note = await getNote(id);
        store.setCurrentOpenTextNode(note);
    };

    return {
        onOpenTreeDialog,
        onCloseTreeDialog,
        onCreateNewFolder,
        onCreateNewNote,
        onDeleteFolder,
        onDropFile,
        onOpenNote,
    };
};

export const useMarkdownTreeNotes = () => {
    const markdownStore = useMarkdownStore();
    const state = useState();
    const actions = useActions(state, markdownStore);

    onBeforeMount(async () => {
        const treeList = await getTreeList();

        markdownStore.setTree(treeList);
    });

    return {
        state,
        actions,
        markdownStore,
    };
};

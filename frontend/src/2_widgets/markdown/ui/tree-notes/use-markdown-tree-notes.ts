import { ref, onBeforeMount } from 'vue';

import { createGroup, createNote, deleteGroup, updateGroup, updateNote } from '@features/markdown/api/';

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

        await createNote({ group_id: String(id), title: 'New Note' });
        onCloseTreeDialog();
    };

    const onDeleteFolder = async (id: string, e: any) => {
        e.preventDefault();
        e.stopPropagation();

        await deleteGroup(id);
        onCloseTreeDialog();
    };

    const onDropFile = async (draggingNode: any, dropNode: any) => {
        if (draggingNode.data.isNote) {
            await updateNote({ id: String(draggingNode.data.id), group_id: String(dropNode.data.id) });
            return;
        }
        await updateGroup({ id: String(draggingNode.data.id), parentGroupId: String(dropNode.data.id) });
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
        const treeListResponse = await getTreeList();

        markdownStore.setTree(treeListResponse);
    });

    return {
        state,
        actions,
        markdownStore,
    };
};

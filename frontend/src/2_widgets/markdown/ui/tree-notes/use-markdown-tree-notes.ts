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

    const updateTreeNotes = async () => {
        const treeListResponse = await getTreeList();

        store.setTree(treeListResponse);
    };

    const onCreateNewFolder = async (id: string, e: any) => {
        e.preventDefault();
        e.stopPropagation();

        await createGroup({ name: 'New Folder', parentIdGroup: String(id) });
        await updateTreeNotes();
    };

    const onEditFolder = async (id: string, currentName: string, e: any) => {
        e.preventDefault();
        e.stopPropagation();

        // eslint-disable-next-line no-alert
        const newValue = prompt(`"${currentName}" переименовать в:`, '');
        await updateGroup({ id: String(id), name: newValue || currentName });
        await updateTreeNotes();
    };

    const onEditNote = async (id: string, currentName: string, e: any) => {
        e.preventDefault();
        e.stopPropagation();
        // eslint-disable-next-line no-alert
        const newValue = prompt(`"${currentName}" переименовать в:`, '');
        await updateNote({ id: String(id), title: newValue || currentName });
        await updateTreeNotes();
    };

    const onCreateNewNote = async (id: string, e: any) => {
        e.preventDefault();
        e.stopPropagation();

        await createNote({ group_id: String(id), title: 'New Note' });
        await updateTreeNotes();
    };

    const onDeleteFolder = async (id: string, e: any) => {
        e.preventDefault();
        e.stopPropagation();

        await deleteGroup(id);
        await updateTreeNotes();
    };

    const onDropFile = async (draggingNode: any, dropNode: any) => {
        if (draggingNode.data.isNote) {
            await updateNote({ id: String(draggingNode.data.id), group_id: String(dropNode.data.id) });
            return;
        }
        await updateGroup({ id: String(draggingNode.data.id), parentGroupId: String(dropNode.data.id) });
        await updateTreeNotes();
    };

    const onOpenNote = async (id: string) => {
        const note = await getNote(id);
        store.setCurrentOpenNodeInfo(note.data);
    };

    return {
        onOpenTreeDialog,
        onCloseTreeDialog,
        onCreateNewFolder,
        onCreateNewNote,
        onDeleteFolder,
        onDropFile,
        onOpenNote,
        onEditFolder,
        onEditNote,
        updateTreeNotes,
    };
};

export const useMarkdownTreeNotes = () => {
    const markdownStore = useMarkdownStore();
    const state = useState();
    const actions = useActions(state, markdownStore);

    onBeforeMount(actions.updateTreeNotes);

    return {
        state,
        actions,
        markdownStore,
    };
};

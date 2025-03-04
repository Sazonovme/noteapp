import { ref } from 'vue';

import { deleteNote as deleteNoteApi } from '@features/markdown/api';

import { useMarkdownStore } from '@entities/markdown';

const useState = () => {
    const dialogVisible = ref(false);

    return {
        dialogVisible,
    };
};

type UseActionsType = ReturnType<typeof useState>;
// eslint-disable-next-line no-unused-vars
const useActions = (emit: (e: string) => void, state: UseActionsType, store: ReturnType<typeof useMarkdownStore>) => {
    const toggleDialog = () => {
        state.dialogVisible.value = !state.dialogVisible.value;
    };

    const deleteNote = async () => {
        state.dialogVisible.value = false;
        await deleteNoteApi(String(store.currentOpenNodeInfo.id));
        location.reload();
        emit('success-delete');
    };

    return {
        deleteNote,
        toggleDialog,
    };
};

export const useDeleteNoteButton = (emit: any) => {
    const markdownStore = useMarkdownStore();
    const state = useState();
    const actions = useActions(emit, state, markdownStore);

    return {
        actions,
        state,
    };
};

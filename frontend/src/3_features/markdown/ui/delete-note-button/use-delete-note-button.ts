import { ref } from 'vue';

const useState = () => {
    const dialogVisible = ref(false);

    return {
        dialogVisible,
    };
};

type UseActionsType = ReturnType<typeof useState>;
// eslint-disable-next-line no-unused-vars
const useActions = (emit: (e: string) => void, state: UseActionsType) => {
    const toggleDialog = () => {
        state.dialogVisible.value = !state.dialogVisible.value;
    };

    const deleteNote = async () => {
        state.dialogVisible.value = false;
        console.log('deleting');
        emit('success-delete');
    };

    return {
        deleteNote,
        toggleDialog,
    };
};

export const useDeleteNoteButton = (emit: any) => {
    const state = useState();
    const actions = useActions(emit, state);

    return {
        actions,
        state,
    };
};

const useActions = (emit: any) => {
    const createNewNote = async () => {
        console.log('creating');

        emit('success-create');
    };

    return {
        createNewNote,
    };
};

export const useCreateNewNoteButton = (emit: any) => {
    const actions = useActions(emit);

    return {
        emit,
        actions,
    };
};

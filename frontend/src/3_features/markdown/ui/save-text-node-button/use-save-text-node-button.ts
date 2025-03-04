import { updateNote } from '@features/markdown/api';

import { useMarkdownStore } from '@entities/markdown';

const useActions = (store: ReturnType<typeof useMarkdownStore>) => {
    const updateNoteText = async () => {
        await updateNote({ id: String(store.currentOpenNodeInfo.id), text: store.currentOpenNodeInfo.text });
    };

    return {
        updateNoteText,
    };
};

export const useSaveTextNodeButton = () => {
    const markdownStore = useMarkdownStore();
    const actions = useActions(markdownStore);

    return {
        markdownStore,
        actions,
    };
};

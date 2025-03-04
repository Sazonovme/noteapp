import { useMarkdownStore } from '@entities/markdown';

export const useNoteTitleChange = () => {
    const markdownStore = useMarkdownStore();

    return {
        markdownStore,
    };
};

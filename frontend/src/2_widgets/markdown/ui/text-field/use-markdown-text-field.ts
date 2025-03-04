import { useMarkdownStore } from '@entities/markdown';

export const useMarkdownTextField = () => {
    const markdownStore = useMarkdownStore();

    return {
        markdownStore,
    };
};

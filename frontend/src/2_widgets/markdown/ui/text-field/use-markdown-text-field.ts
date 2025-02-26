import { onBeforeMount, ref } from 'vue';

import { useMarkdownStore } from '@entities/markdown';
import { getNote } from '@entities/markdown/api';

const MSG = `
# Привет я заголовок
## поменьше
__жирдяй__
~~черкаш~~
- [x] таска ок
- [ ] таска не ок
- обычный текст
`;

const useState = () => {
    const noteText = ref<string>(MSG);

    return {
        noteText,
    };
};

export const useMarkdownTextField = () => {
    const markdownStore = useMarkdownStore();
    const state = useState();

    onBeforeMount(async () => {
        // state.noteText.value = 'Hello world!';
        // state.noteText.value = await getNote(markdownStore.currentOpenIdNode);
    });

    return {
        state,
        markdownStore,
    };
};

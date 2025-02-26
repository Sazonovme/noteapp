import { defineStore } from 'pinia';
import { ref } from 'vue';

import { storage } from '@shared/utils';

import { MARKDOWN_MODE_LOCAL } from '../constants';

export const useMarkdownStore = defineStore('markdown', () => {
    const readOnly = ref((storage.local.getItem(MARKDOWN_MODE_LOCAL) ?? 'true') === 'true');

    const setReadOnly = (val: boolean) => {
        readOnly.value = val;
        storage.local.setItem(MARKDOWN_MODE_LOCAL, String(val));
    };

    const tree = ref<Record<string, any>[]>([]);
    const setTree = (val: {}[]) => {
        tree.value = val;
    };

    const currentOpenIdNode = ref('');
    const setCurrentOpenIdNode = (val: string) => {
        currentOpenIdNode.value = val;
    };

    const currentOpenTextNode = ref('');
    const setCurrentOpenTextNode = (val: string) => {
        currentOpenTextNode.value = val;
    };

    return {
        readOnly,
        setReadOnly,

        tree,
        setTree,

        currentOpenIdNode,
        setCurrentOpenIdNode,

        currentOpenTextNode,
        setCurrentOpenTextNode,
    };
});

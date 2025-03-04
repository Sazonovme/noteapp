import { defineStore } from 'pinia';
import { ref } from 'vue';

import { storage } from '@shared/utils';

import { MARKDOWN_MODE_LOCAL } from '../constants';
import { NodeInfoType } from '../model';

const DEFAULT_NODE_INFO: NodeInfoType = {
    group_id: 0,
    id: 0,
    text: 'Default text Node',
    title: 'Default title Node',
    user_email: 'Root',
};

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

    const currentOpenNodeInfo = ref(DEFAULT_NODE_INFO);
    const setCurrentOpenNodeInfo = (val: NodeInfoType) => {
        currentOpenNodeInfo.value = val;
    };

    return {
        readOnly,
        setReadOnly,

        tree,
        setTree,

        currentOpenNodeInfo,
        setCurrentOpenNodeInfo,
    };
});

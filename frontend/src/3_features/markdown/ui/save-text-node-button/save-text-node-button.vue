<template>
    <ElButton
        :class="'delete-note-button ' + $attrs.class"
        type="info"
        :size="props.size || 'small'"
        :circle="props.circle"
        plain
        @click="actions.updateNoteText"
    >
        <slot>
            Сохранить
        </slot>
    </ElButton>
</template>

<script setup lang="ts">
import { ElButton } from 'element-plus';

import { useSaveTextNodeButton } from './use-save-text-node-button';

const props = defineProps(['circle', 'size']);

const { actions } = useSaveTextNodeButton();
</script>

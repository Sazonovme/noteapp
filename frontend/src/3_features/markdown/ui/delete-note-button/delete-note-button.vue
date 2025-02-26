<template>
    <div>
        <ElButton
            :class="'delete-note-button ' + $attrs.class"
            type="danger"
            :size="props.size || 'small'"
            :circle="props.circle"
            plain
            @click="actions.toggleDialog"
        >
            <slot>
                Удалить
            </slot>
        </ElButton>
        <ElDialog
            v-model="state.dialogVisible.value"
            title="Вы уверены?"
            width="500"
        >
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="actions.toggleDialog">
                        Отмена
                    </el-button>
                    <el-button
                        type="danger"
                        @click="actions.deleteNote"
                    >
                        Удалить
                    </el-button>
                </div>
            </template>
        </ElDialog>
    </div>
</template>

<script setup lang="ts">
import { ElButton, ElDialog } from 'element-plus';

import { useDeleteNoteButton } from './use-delete-note-button';

const props = defineProps(['circle', 'size']);
const emit = defineEmits(['success-delete']);

const { actions, state } = useDeleteNoteButton(emit);
</script>

<style scoped src="./delete-note-button.css" />

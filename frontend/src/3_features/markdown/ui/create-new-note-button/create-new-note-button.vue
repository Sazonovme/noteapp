<template>
    <ElButton
        type="primary"
        size="large"
        :class="'create-new-note-button ' + $attrs.class"
        @click="actions.createNewNote"
    >
        Создать новую заметку
    </ElButton>
</template>

<script setup lang="ts">
import { ElButton } from 'element-plus';

import { useCreateNewNoteButton } from './use-create-new-note-button';

const emit = defineEmits(['success-create']);

const { actions } = useCreateNewNoteButton(emit);
</script>

<style scoped src="./create-new-note-button.css" />

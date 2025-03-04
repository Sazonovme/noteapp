<template>
    <ElInput
        v-model="markdownStore.currentOpenNodeInfo.title"
        :class="'input ' + $attrs.class"
        size="large"
        disabled
    />
</template>

<script setup lang="ts">
import { ElInput } from 'element-plus';

import { useNoteTitleChange } from './use-note-title-change';

const { markdownStore } = useNoteTitleChange();
</script>

<style src="./note-title-change.css" scoped />

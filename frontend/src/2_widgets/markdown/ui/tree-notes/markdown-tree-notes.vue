<template>
    <button
        class="button"
        @click="setOpen"
    />
    <ElDrawer
        v-model="isOpen"
        title="Опачки! не нужный title"
        direction="ltr"
    >
        <span>Тут будет список заметок!</span>
    </ElDrawer>
</template>

<script setup lang="ts">
import { ElDrawer } from 'element-plus';
import { ref } from 'vue';

import { useMarkdownStore } from '@entities/markdown';

const markdownStore = useMarkdownStore();

const isOpen = ref(false);

const setOpen = () => {
    isOpen.value = true;
};
</script>

<style src="./markdown-tree-notes.css" scoped />

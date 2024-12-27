<template>
    <button
        class="button"
        @click="toggleOpen"
    />
    <ElDrawer
        v-model="isOpen"
        title="Опачки! не нужный title"
        direction="ltr"
    >
        <div>
            <CreateNewNoteButton @success-create="toggleOpen" />
        </div>
    </ElDrawer>
</template>

<script setup lang="ts">
import { ElDrawer } from 'element-plus';
import { ref } from 'vue';

import { CreateNewNoteButton } from '@features/markdown';

// import { useMarkdownStore } from '@entities/markdown';
// const markdownStore = useMarkdownStore();

const isOpen = ref(false);

const toggleOpen = () => {
    isOpen.value = !isOpen.value;
};
</script>

<style src="./markdown-tree-notes.css" scoped />

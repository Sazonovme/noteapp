<template>
    <div class="text-container">
        <MarkdownRender
            v-if="markdownStore.readOnly"
            :markdown-text="text"
        />
        <MarkdownWriteArea
            v-else
            v-model="text"
        />
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { MarkdownWriteArea } from '@features/markdown';

import { MarkdownRender, useMarkdownStore } from '@entities/markdown';

const markdownStore = useMarkdownStore();

const MSG = `
# markdown-it rulezz!

__13d__ 
~~dslfkjn~~ //\n
- [x] das;k
- asdk2
`;

const text = ref(MSG);
</script>

<style src="./markdown-text-field.css" scoped />

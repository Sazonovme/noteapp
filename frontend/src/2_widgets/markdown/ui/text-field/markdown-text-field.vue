<template>
    <div class="text-container">
        <MarkdownRender
            v-if="markdownStore.readOnly"
            :markdown-text="markdownStore.currentOpenTextNode"
        />
        <MarkdownWriteArea
            v-else
            v-model="markdownStore.currentOpenTextNode"
        />
    </div>
</template>

<script setup lang="ts">
import { MarkdownWriteArea } from '@features/markdown';

import { MarkdownRender } from '@entities/markdown';

import { useMarkdownTextField } from './use-markdown-text-field';

const { state, markdownStore } = useMarkdownTextField();
</script>

<style src="./markdown-text-field.css" scoped />

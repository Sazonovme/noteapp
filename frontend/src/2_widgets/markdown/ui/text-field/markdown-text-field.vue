<template>
    <div>
        <NoteTitleChange />

        <div class="text-container">
            <MarkdownRender
                v-if="markdownStore.readOnly"
                :markdown-text="markdownStore.currentOpenNodeInfo.text"
            />
            <MarkdownWriteArea
                v-else
                v-model="markdownStore.currentOpenNodeInfo.text"
            />
        </div>
    </div>
</template>

<script setup lang="ts">
import { MarkdownWriteArea, NoteTitleChange } from '@features/markdown';

import { MarkdownRender } from '@entities/markdown';

import { useMarkdownTextField } from './use-markdown-text-field';

const { markdownStore } = useMarkdownTextField();
</script>

<style src="./markdown-text-field.css" scoped />

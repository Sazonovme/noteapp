<template>
    <MarkdownTreeNotes />
    <div class="container">
        <div
            v-if="markdownStore.currentOpenNodeInfo.id !== 0"
            class="btn-container"
        >
            <MarkdownToggleReadOnly class="switcher" />

            <SaveTextNodeButton
                size="large"
                circle
            >
                <div class="save-img" />
            </SaveTextNodeButton>
            <DeleteNoteButton
                size="large"
                circle
            >
                <div class="delete-img" />
            </DeleteNoteButton>
        </div>

        <MarkdownWidget />
    </div>
</template>

<script setup lang="ts">
import { MarkdownWidget, MarkdownTreeNotes } from '@widgets/markdown';

import { MarkdownToggleReadOnly, DeleteNoteButton, SaveTextNodeButton } from '@features/markdown';

import { useMarkdownStore } from '@entities/markdown';

const markdownStore = useMarkdownStore();
</script>

<style src="./page-home.css" scoped />

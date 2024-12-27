<template>
    <MarkdownTreeNotes />
    <div class="container">
        <MarkdownToggleReadOnly class="switcher" />
        <DeleteNoteButton
            class="delete-btn"
            size="large"
            circle
        >
            <div class="img" />
        </DeleteNoteButton>

        <MarkdownWidget />
    </div>
</template>

<script setup lang="ts">
import { MarkdownWidget, MarkdownTreeNotes } from '@widgets/markdown';

import { MarkdownToggleReadOnly, DeleteNoteButton } from '@features/markdown';

</script>

<style src="./page-home.css" scoped />

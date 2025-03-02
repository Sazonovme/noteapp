<template>
    <button
        class="button"
        @click="actions.onOpenTreeDialog"
    >
        <div class="button-img" />
    </button>
    <ElDrawer
        v-model="state.isTreeDialogOpen.value"
        title="Крутое меню"
        direction="ltr"
    >
        <div>
            <ElTree
                :data="markdownStore.tree"
                :allow-drop="(_: any, dropNode: any, dropType: string) => (dropType === 'inner') && dropNode.data.isFolder"
                draggable
                @node-drop="actions.onDropFile"
            >
                <template #default="{ data }">
                    <span
                        class="custom-tree-node"
                        @click="() => data.isNote && actions.onOpenNote(data.id)"
                    >
                        <span @click="() => {console.log(data)}">{{ data.title }}</span>
                        <span class="buttons">
                            <div
                                v-if="data.isFolder"
                                class="img img-document-plus"
                                @click="(e: any) => actions.onCreateNewNote(data.id, e)"
                            />
                            <div
                                v-if="data.isFolder"
                                class="img img-folder-plus"
                                @click="(e: any) => actions.onCreateNewFolder(data.id, e)"
                            />
                            <div
                                v-if="data.isFolder"
                                class="img img-folder-minus"
                                @click="(e: any) => actions.onDeleteFolder(data.id, e)"
                            />
                        </span>
                    </span>
                </template>
            </ElTree>
            <!-- <CreateNewNoteButton @success-create="onCloseTreeDialog" /> -->
        </div>
    </ElDrawer>
</template>

<script setup lang="ts">
import { ElDrawer, ElTree } from 'element-plus';

import { useMarkdownTreeNotes } from './use-markdown-tree-notes';

const { state, actions, markdownStore } = useMarkdownTreeNotes();
</script>

<style src="./markdown-tree-notes.css" scoped />

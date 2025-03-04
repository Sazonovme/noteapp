<template>
    <div>
        <ElSwitch
            size="default"
            :model-value="markdownStore.readOnly"
            @change="markdownStore.setReadOnly"
        />
    </div>
</template>

<script setup lang="ts">
import { ElSwitch } from 'element-plus';

import { useMarkdownStore } from '@entities/markdown';

const markdownStore = useMarkdownStore();
</script>

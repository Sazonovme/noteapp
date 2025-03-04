<template>
    <ElButton
        :class="$attrs.class"
        type="info"
        size="large"
        circle
        plain
        @click="props.to ? $router.push(props.to) : $router.go(-1)"
    >
        <slot>
            Назад
        </slot>
    </ElButton>
</template>

<script setup lang="ts">
import { ElButton } from 'element-plus';

const props = defineProps(['to']);
</script>

<template>
    <RouterLink :to="{ name: 'auth', query: { type: 'signin' } }">
        <ElButton
            type="primary"
            size="large"
            circle
            plain
        >
            <div class="logout-img" />
        </ElButton>
    </RouterLink>
</template>

<script setup lang="ts">
import { ElButton } from 'element-plus';
</script>

<style src="./logout-button.css" scoped />

<template>
    <RouterLink :to="{ name: 'auth', query: { type: 'signin' } }">
        <ElButton
            type="primary"
            size="large"
            circle
            plain
        >
            <div class="user-img" />
        </ElButton>
    </RouterLink>
</template>

<script setup lang="ts">
import { ElButton } from 'element-plus';
</script>

<style src="./sign-in-button.css" scoped />

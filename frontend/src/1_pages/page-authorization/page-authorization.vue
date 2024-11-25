<template>
    <div
        ref="backgroundRef"
        class="background"
    />
    <div class="container">
        Форма авторизации
    </div>
</template>

<script setup lang="ts">
import { usePageAuthorization } from './use-page-authorization';

const { backgroundRef } = usePageAuthorization();
</script>

<style scoped src="./page-authorization.css" />

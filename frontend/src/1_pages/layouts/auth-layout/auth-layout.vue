<template>
    <BackButton
        class="link"
        :to="ROUTES_PATH_COMMON.HOME"
    >
        На главную
    </BackButton>
    <div class="container">
        <main class="main">
            <slot />
        </main>
    </div>
</template>

<script setup lang="ts">
import { BackButton } from '@shared/ui';
import { ROUTES_PATH_COMMON } from '@shared/constants';
</script>

<style src="./auth-layout.css" scoped />

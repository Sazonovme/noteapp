<template>
    <div class="container">
        <!-- todo это нужно вынести в отельный компонент -->
        <nav>
            <SignInButton class="auth-button" />
        </nav>

        <main>
            <slot />
        </main>
    </div>
</template>

<script setup lang="ts">
import { SignInButton } from '@features/user';
</script>

<style src="./base-layout.css" scoped />

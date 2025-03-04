<template>
    <div class="container">
        <nav>
            <LogoutButton
                v-if="isAuth"
                class="auth-button"
            />
            <SignInButton
                v-else
                class="auth-button"
            />
        </nav>

        <main>
            <slot />
        </main>
    </div>
</template>

<script setup lang="ts">
import { SignInButton, LogoutButton } from '@features/user';

import { useCheckIsAuth } from '@entities/user';

const { isAuth } = useCheckIsAuth();
</script>

<style src="./base-layout.css" scoped />

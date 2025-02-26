<template>
    <div
        ref="backgroundRef"
        class="background"
    />
    <div class="container">
        <div class="substrate" />
        <div class="forms">
            <SignIn v-if="$route.query.type === 'signin'" />
            <SignUp v-if="$route.query.type === 'signup'" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { SignIn, SignUp } from './components';
import { usePageAuthorization } from './use-page-authorization';

const { backgroundRef } = usePageAuthorization();
</script>

<style scoped src="./page-authorization.css" />

<template>
    <!-- todo если это реально рабочий вариант, то кринж -->
    <BaseLayout v-if="$route.meta.layout == 'base'">
        <RouterView />
    </BaseLayout>
    <AuthLayout v-else-if="$route.meta.layout == 'auth'">
        <RouterView />
    </AuthLayout>
    <RouterView v-else />
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router';

import { BaseLayout, AuthLayout } from '@pages/layouts';
</script>

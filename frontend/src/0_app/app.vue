<template>
    <!-- todo если это реально рабочий вариант, то кринж -->
    <BaseLayout v-if="$route.meta.layout == 'base'">
        <RouterView />
    </BaseLayout>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router';

import { BaseLayout } from 'pages/layouts/base-layout';
</script>

<template>
    <div class="container">
        <!-- todo это нужно вынести в отельный компонент -->
        <nav>
            <RouterLink to="/">
                Home
            </RouterLink>
            <RouterLink to="/auth">
                auth
            </RouterLink>
        </nav>

        <main>
            <slot />
        </main>
    </div>
</template>

<script setup lang="ts">
</script>

<style scoped src="./base-layout.css" />

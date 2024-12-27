<template>
    <button
        :class="'text ' + $attrs.class"
        @click="props.to ? $router.push(props.to) : $router.go(-1)"
    >
        <slot>
            Назад
        </slot>
    </button>
</template>

<script setup lang="ts">
const props = defineProps(['to']);
</script>

<style src="./back-button.css" scoped />

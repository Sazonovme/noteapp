<template>
    <textarea
        class="textarea"
        :value="props.modelValue"
        @input="$event => emit('update:modelValue', ($event.target as unknown as HTMLInputElement)?.value || '')"
    />
</template>

<script setup lang="ts">
const props = defineProps(['modelValue']);
const emit = defineEmits(['update:modelValue']);
</script>

<style src="./markdown-write-area.css" scoped />

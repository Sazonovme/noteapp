<template>
    <div v-html="md.render(props.markdownText)" />
</template>

<script setup lang="ts">
import MarkdownIt from 'markdown-it';
import MarkdownItAbbr from 'markdown-it-abbr';
import MarkdownItAnchor from 'markdown-it-anchor';
import MarkdownItFootnote from 'markdown-it-footnote';
import MarkdownItHighlightjs from 'markdown-it-highlightjs';
import MarkdownItSub from 'markdown-it-sub';
import MarkdownItSup from 'markdown-it-sup';
import MarkdownItTasklists from 'markdown-it-task-lists';

const props = defineProps(['markdownText']);

const md = new MarkdownIt({
    html: true, // разрешить HTML внутри Markdown
    linkify: true, // преобразование URL в кликабельные ссылки
    typographer: true, // улучшенная типографика
})
    .use(MarkdownItAbbr)
    .use(MarkdownItAnchor)
    .use(MarkdownItFootnote)
    .use(MarkdownItHighlightjs)
    .use(MarkdownItSub)
    .use(MarkdownItSup)
    .use(MarkdownItTasklists);
// .use(MarkdownItTOC);
</script>

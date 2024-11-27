import path from 'path';

import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

// https://vite.dev/config/
export default defineConfig({
    plugins: [vue()],
    server: {
        port: 8080,
        open: true,
    },
    resolve: {
        alias: {
            '@app': path.resolve(__dirname, './src/0_app'),
            '@pages': path.resolve(__dirname, './src/1_pages'),
            '@widgets': path.resolve(__dirname, './src/2_widgets'),
            '@features': path.resolve(__dirname, './src/3_features'),
            '@entities': path.resolve(__dirname, './src/4_entities'),
            '@shared': path.resolve(__dirname, './src/shared'),
            '@': path.resolve(__dirname, './src'),
        },
        extensions: ['.ts', '.tsx', '.js', 'jsx', '.mjs', '.css', 'd.ts', '.vue'],
    },
});

const path = require('path');

module.exports = {
    env: {
        browser: true,
        node: true,
        amd: true,
        es6: true,
        es2017: true,
        es2021: true,
    },
    parser: 'vue-eslint-parser',
    parserOptions: {
        parser: '@typescript-eslint/parser',
        // tsconfigRootDir: __dirname,
        // project: [
        //     path.resolve(__dirname, './tsconfig.app.json'),
        // ],
        extraFileExtensions: ['.vue', '.json'],
    },
    extends: [
        'eslint:recommended',
        // 'plugin:@typescript-eslint/recommended',
        'plugin:vue/vue3-recommended',
        './rules/eslint/index.cjs',
    ],
    plugins: [
        // '@typescript-eslint',
        // 'vue',
        'import',
    ],
    rules: {
        'generator-star-spacing': 'off',
        // allow debugger during development
        'no-debugger': 'error',
        indent: ['error', 4],
        'vue/html-indent': ['error', 4],
        'vue/script-indent': ['error', 4],
    },
    // overrides: [
    //     {
    //         files: [
    //             '**/*.vue'
    //         ],
    //         parser: 'vue-eslint-parser',
    //         parserOptions: {
    //             parser: '@typescript-eslint/parser',
    //             tsconfigRootDir: __dirname,
    //             // project: [
    //             //     path.resolve(__dirname, './tsconfig.app.json'),
    //             // ],
    //             // extraFileExtensions: ['.vue'],
    //         },
    //         extends: [
    //             'eslint:recommended',
    //             // "plugin:@typescript-eslint/eslint-recommended",
    //             // "plugin:@typescript-eslint/recommended",
    //             // "plugin:vue/vue3-recommended",

    //             // 'plugin:@typescript-eslint/recommended',
    //             'plugin:vue/vue3-recommended',
    //             // "@vue/eslint-config-typescript",
    //             './rules/eslint/index.cjs',
    //         ],
    //         plugins: [
    //             // '@typescript-eslint',
    //             // '@typescript-eslint/eslint-plugin',
    //             // 'plugin:vue/vue3-recommended',
    //             'vue',
    //             'import',
    //         ],
            
    //         rules: {
    //             // allow async-await
    //             'generator-star-spacing': 'off',
    //             // allow debugger during development
    //             'no-debugger': 'error',
    //             indent: ['error', 4],
    //             'vue/html-indent': ['error', 4],
    //             'vue/script-indent': ['error', 4],
    //         },
    //     },
    // ],
};

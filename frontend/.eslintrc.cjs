const path = require('path');

module.exports = {
    globals: {
        NodeJS: true,
    },
    env: {
        browser: true,
        node: true,
        amd: true,
        es6: true,
        es2017: true,
        es2021: true,
    },
    // parser: '@typescript-eslint/parser',
    parser: 'vue-eslint-parser',
    parserOptions: {
        // parser: 'vue-eslint-parser',
        parser: '@typescript-eslint/parser',
        sourceType: 'module',
        tsconfigRootDir: __dirname,
        project: [
            // path.resolve(__dirname, './tsconfig.json'),
            path.resolve(__dirname, './tsconfig.app.json'),
        ],
        ecmaVersion: 'latest',
        ecmaFeatures: {
            jsx: true,
        },
        extraFileExtensions: ['.vue'],
    },
    extends: ['eslint:recommended', 'plugin:@typescript-eslint/recommended', 'plugin:vue/vue3-recommended', './rules/eslint/index.cjs'/* , 'plugin:security/recommended' */],
    plugins: ['@typescript-eslint', 'vue', 'import'],
    rules: {
        // allow async-await
        'generator-star-spacing': 'off',
        // allow debugger during development
        'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
        indent: ['error', 4],
        'vue/html-indent': ['error', 4],
        'vue/script-indent': ['error', 4],
    },
    root: true,
};

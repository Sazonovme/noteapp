module.exports = {
    settings: {
        'import/resolver': {
            node: {
                extensions: ['.js', '.cjs', '.jsx', '.ts', '.tsx', '.json'],
            },
        },
        'import/extensions': ['.js', '.cjs', '.jsx', '.ts', '.tsx', '.json'],
        'import/core-modules': [],
        'import/external-module-folders': ['node_modules'],
        'import/ignore': [
            'node_modules',
            '\\.(coffee|scss|css|less|hbs|svg|json|png|jpg|gif)$',
        ],
    },
    rules: {
        'import/no-unresolved': ['error', {
            commonjs: true,
            caseSensitive: true,
            ignore: ['shared/*', 'src/*', 'pages/*', 'app/*', 'features/*', 'widgets/*', 'entities/*'],
        }],
        'import/named': 'off',
        'import/default': 'off',
        'import/namespace': 'off',
        'import/export': 'error',
        'import/no-named-as-default': 'warn',
        'import/no-named-as-default-member': 'warn',
        'import/no-deprecated': 'off',
        'import/no-extraneous-dependencies': 'off',
        'import/no-mutable-exports': 'error',
        'import/no-commonjs': 'off',
        'import/no-amd': 'error',
// TODO: enable?
        'import/no-nodejs-modules': 'off',
        'import/first': ['warn', 'absolute-first'],
        'import/imports-first': 'off',
        'import/no-duplicates': ['error', { 'prefer-inline': true }],
// TODO: enable?
        'import/no-namespace': 'off',
        'import/extensions': ['off', 'always', {
            js: 'never',
            jsx: 'never',
            json: 'always',
            ts: 'never',
            tsx: 'never',
        }],
        'import/order': ['error', {
            groups: ['builtin', 'external', 'unknown', 'internal', 'parent', 'sibling', 'index'],
            pathGroups: [
                {
                    pattern: 'app/**',
                    group: 'unknown',
                    position: 'after',
                },
                {
                    pattern: '{src/1_,}pages{/**,}',
                    group: 'unknown',
                    position: 'after',
                },
                {
                    pattern: '{src/2_,}widgets{/**,}',
                    group: 'unknown',
                    position: 'after',
                },
                {
                    pattern: '{src/3_,}features{/**,}',
                    group: 'unknown',
                    position: 'after',
                },
                {
                    pattern: '{src/4_,}entities{/**,}',
                    group: 'unknown',
                    position: 'after',
                },
                {
                    pattern: 'src/**',
                    group: 'unknown',
                    position: 'after',
                },
                {
                    pattern: 'shared/**',
                    group: 'unknown',
                    position: 'after',
                },
            ],
            pathGroupsExcludedImportTypes: ['builtin'],
            // distinctGroup: false,
            'newlines-between': 'always',
        }],
        'import/newline-after-import': 'warn',
        'import/prefer-default-export': 'off',
        'import/no-restricted-paths': 'off',
        'import/max-dependencies': ['off', { max: 10 }],
        'import/no-absolute-path': 'error',
        'import/no-dynamic-require': 'error',
        'import/no-internal-modules': ['off', {
            allow: [],
        }],
        'import/unambiguous': 'off',
        'import/no-webpack-loader-syntax': 'warn',
        'import/no-unassigned-import': 'off',
        'import/no-named-default': 'warn',
    },
};

module.exports = {
    rules: {
        'array-bracket-newline': ['error', 'consistent'], // перенос строк в массиве
        'array-bracket-spacing': ['error', 'never'], // запретить пробелы внутри скобок
        'array-element-newline': 'off',
        'block-spacing': ['error', 'always'],
        'brace-style': ['error', '1tbs', { allowSingleLine: true }],
        camelcase: ['error', { properties: 'never' }],
        'capitalized-comments': 'off',
        'comma-dangle': [
            'error',
            {
                arrays: 'always-multiline',
                objects: 'always-multiline',
                imports: 'never',
                exports: 'always-multiline',
                functions: 'never',
            },
        ],
        'comma-spacing': ['error', { before: false, after: true }],
        'comma-style': [
            'error',
            'last',
            {
                exceptions: {
                    ArrayExpression: false,
                    ArrayPattern: false,
                    ArrowFunctionExpression: false,
                    CallExpression: false,
                    FunctionDeclaration: false,
                    FunctionExpression: false,
                    ImportDeclaration: false,
                    ObjectExpression: false,
                    ObjectPattern: false,
                    VariableDeclaration: false,
                    NewExpression: false,
                },
            },
        ],
        'computed-property-spacing': ['error', 'never'],
        'consistent-this': 'off',
        'eol-last': ['error', 'always'],
        'func-call-spacing': ['error', 'never'],
        'func-name-matching': 'off',
        'func-names': 'warn',
        'func-style': 'off',
        'function-call-argument-newline': ['error', 'consistent'],
        'function-paren-newline': ['error', 'consistent'],
        'id-denylist': 'off',
        'id-length': 'off',
        'id-match': 'off',
        'implicit-arrow-linebreak': ['error', 'beside'],
        'jsx-quotes': 'off',
        'key-spacing': ['error', { beforeColon: false, afterColon: true }],
        'keyword-spacing': [
            'error',
            {
                before: true,
                after: true,
            },
        ],
        'line-comment-position': 'off',
        'linebreak-style': 'off',
        'lines-around-comment': 'off',
        'lines-between-class-members': ['error', 'always', { exceptAfterSingleLine: false }],
        'max-depth': 'off',
        'max-len': [
            //max length line
            'off',
            999,
            4,
            {
                ignoreUrls: true,
                ignoreComments: true,
                ignoreRegExpLiterals: true,
                ignoreStrings: true,
                ignoreTemplateLiterals: true,
            },
        ],
        'max-lines': 'off',
        'max-lines-per-function': 'off',
        'max-nested-callbacks': 'off',
        'max-params': 'off',
        'max-statements': 'off',
        'max-statements-per-line': 'off',
        'multiline-comment-style': 'off',
        'multiline-ternary': 'off',
        'new-cap': [
            'error',
            {
                newIsCap: true,
                capIsNew: false,
            },
        ],
        'new-parens': 'error',
        'newline-per-chained-call': ['error', { ignoreChainWithDepth: 4 }],
        'no-array-constructor': 'error',
        'no-bitwise': 'error',
        'no-continue': 'off',
        'no-inline-comments': 'off',
        'no-lonely-if': 'off',
        'no-mixed-spaces-and-tabs': 'error',
        'no-multi-assign': 'error',
        'no-multiple-empty-lines': ['error', { max: 1, maxBOF: 0, maxEOF: 0 }],
        'no-negated-condition': 'off',
        'no-nested-ternary': 'off',
        'no-new-object': 'error',
        'no-plusplus': 'off',
        'no-restricted-syntax': 'off',
        'no-ternary': 'off',
        'no-tabs': 'error',
        'no-trailing-spaces': 'error',
        'no-underscore-dangle': 'off',
        'no-unneeded-ternary': 'error',
        'no-whitespace-before-property': 'error',
        'nonblock-statement-body-position': 'off',
        'object-curly-newline': [
            'error',
            {
                ObjectExpression: { minProperties: 10, multiline: true, consistent: true },
                ObjectPattern: { minProperties: 10, multiline: true, consistent: true },
                ImportDeclaration: { minProperties: 10, multiline: true, consistent: true },
                ExportDeclaration: { minProperties: 10, multiline: true, consistent: true },
            },
        ],
        'object-curly-spacing': ['error', 'always'],
        'object-property-newline': [
            'error',
            {
                allowAllPropertiesOnSameLine: true,
            },
        ],
        'one-var': 'off',
        'one-var-declaration-per-line': ['error', 'always'],
        'operator-assignment': ['error', 'always'],
        // 'operator-linebreak': ['error', 'before'], // Говрить как переносить ? : и прочее
        'padded-blocks': [
            'error',
            {
                blocks: 'never',
                classes: 'never',
                switches: 'never',
            },
        ],
        'padding-line-between-statements': 'off',
        'prefer-exponentiation-operator': 'off',
        'prefer-object-spread': 'off',
        'quote-props': ['error', 'as-needed', { keywords: false, unnecessary: true, numbers: false }],
        quotes: ['error', 'single', { avoidEscape: true }],
        semi: ['error', 'always'],
        'semi-spacing': ['error', { before: false, after: true }],
        'semi-style': ['error', 'last'],
        'sort-keys': 'off',
        'sort-vars': 'off',
        'space-before-blocks': 'error',
        'space-before-function-paren': [
            'error',
            {
                anonymous: 'never',
                named: 'never',
                asyncArrow: 'always',
            },
        ],
        'space-in-parens': ['error', 'never'],
        'space-infix-ops': 'error',
        'space-unary-ops': ['error', { words: true, nonwords: false }],
        'spaced-comment': 'off',
        'switch-colon-spacing': ['error', { after: true, before: false }],
        'template-tag-spacing': ['error', 'never'],
        'unicode-bom': 'off',
        'wrap-regex': 'off',
    },
};

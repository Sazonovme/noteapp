module.exports = {
    rules: {
        'color-no-invalid-hex': true,
        'font-family-no-duplicate-names': true,
        'font-family-no-missing-generic-family-keyword': true,
        'function-calc-no-unspaced-operator': true,
        'string-no-newline': true,
        'unit-no-unknown': true,
        'property-no-unknown': true,
        'keyframe-declaration-no-important': true,
        'declaration-block-no-duplicate-properties': true,
        'block-no-empty': true,
        'selector-pseudo-class-no-unknown': true,
        'selector-pseudo-element-no-unknown': true,
        'selector-type-no-unknown': true,
        'media-feature-name-no-unknown': true,
        'at-rule-no-unknown': true,
        'no-descending-specificity': true,
        'no-duplicate-at-import-rules': true,
        'no-duplicate-selectors': true,
        'no-empty-source': true,
        'no-extra-semicolons': true,
        'no-invalid-double-slash-comments': true,
        // 'color-no-hex': true,
        'length-zero-no-unit': true,
        'font-weight-notation': 'numeric',
        'number-max-precision': 2,
        'unit-disallowed-list': [
            ['rem', 'em'],
            {
                ignoreMediaFeatureNames: {
                    px: ['min-width', 'max-width'],
                },
            },
        ],
        'shorthand-property-no-redundant-values': true,
        'value-no-vendor-prefix': true,
        // 'property-no-vendor-prefix': true,
        'declaration-block-no-redundant-longhand-properties': true,
        'declaration-block-single-line-max-declarations': 1,
        'selector-max-class': 4,
        'selector-max-id': 1,
        'selector-max-type': 4,
        'selector-no-vendor-prefix': true,
        'selector-pseudo-element-colon-notation': 'single',
        'at-rule-no-vendor-prefix': true,
        'max-nesting-depth': 3,
        'color-hex-case': 'lower',
        'color-hex-length': 'long',
        'font-family-name-quotes': 'always-where-required',
        'function-comma-space-after': 'always',
        'function-comma-space-before': 'never',
        'function-max-empty-lines': 0,
        'function-name-case': 'lower',
        'function-parentheses-space-inside': 'never',
        'function-url-quotes': 'always',
        'number-no-trailing-zeros': true,
        'string-quotes': 'single',
        'unit-case': 'lower',
        'value-keyword-case': 'lower',
        'value-list-comma-space-after': 'always-single-line',
        'value-list-comma-space-before': 'never',
        'value-list-max-empty-lines': 0,
        'property-case': 'lower',
        'declaration-bang-space-after': 'never',
        'declaration-bang-space-before': 'always',
        'declaration-colon-space-after': 'always',
        'declaration-colon-space-before': 'never',
        'declaration-block-semicolon-newline-after': 'always',
        'declaration-block-trailing-semicolon': 'always',
        'block-closing-brace-empty-line-before': 'never',
        'block-closing-brace-newline-after': 'always',
        'block-closing-brace-newline-before': 'always',
        'block-opening-brace-newline-after': 'always',
        'block-opening-brace-space-before': 'always',
        'selector-attribute-brackets-space-inside': 'never',
        'selector-attribute-operator-space-after': 'never',
        'selector-attribute-operator-space-before': 'never',
        'selector-combinator-space-after': 'always',
        'selector-combinator-space-before': 'always',
        'selector-pseudo-class-case': 'lower',
        'selector-pseudo-class-parentheses-space-inside': 'never',
        'selector-pseudo-element-case': 'lower',
        'selector-type-case': 'lower',
        'selector-list-comma-newline-after': 'always',
        'selector-list-comma-newline-before': 'never-multi-line',
        'selector-list-comma-space-before': 'never',
        'rule-empty-line-before': ['always', { except: ['inside-block'] }],
        'media-feature-colon-space-after': 'always',
        'media-feature-colon-space-before': 'never',
        'media-feature-name-case': 'lower',
        'media-feature-parentheses-space-inside': 'never',
        'media-feature-range-operator-space-after': 'always',
        'media-feature-range-operator-space-before': 'always',
        'media-query-list-comma-space-after': 'always',
        'media-query-list-comma-space-before': 'never',
        'at-rule-empty-line-before': 'always',
        'at-rule-name-case': 'lower',
        'at-rule-name-space-after': 'always',
        'at-rule-semicolon-newline-after': 'always',
        'at-rule-semicolon-space-before': 'never',
        indentation: 4,
        'max-empty-lines': 1,
        'no-empty-first-line': true,
        'plugin/declaration-block-no-ignored-properties': true,
        // 'plugin/rational-order': [
        //     true,
        //     {
        //         'empty-line-between-groups': true,
        //     },
        // ],
    },
};
module.exports = {
    rules: {
        'init-declarations': 'off',
        'no-delete-var': 'off',
        'no-label-var': 'error',
        'no-restricted-globals': 'off', // запретить определенные глобальные объекты
        'no-shadow': 'off', // запретить объявление переменных, уже объявленных во внешней области видимости (баг изза typescript-eslint/no-shadow)
        'no-shadow-restricted-names': 'error',
        'no-undef': 'error', // запретить использование необъявленных переменных, если они не упомянуты в global
        'no-undef-init': 'error', // запретить использование undefined при инициализации переменных
        'no-undefined': 'off', // запретить использование undefined
        // запретить объявление переменных, которые не используются в коде
        'no-unused-vars': ['error', { vars: 'all', args: 'after-used', ignoreRestSiblings: true }],
        'no-use-before-define': 'error', // запретить использование переменных до их определения
    },
};

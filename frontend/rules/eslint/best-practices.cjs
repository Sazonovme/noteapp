module.exports = {
    rules: {
        'accessor-pairs': 'off', // принудительно применяет пары геттер / сеттер в объектах
        'array-callback-return': ['warn'], // принудительное выполнение операторов возврата в обратных вызовах методов массива
        'block-scoped-var': 'error',
        'class-methods-use-this': [
            'warn',
            {
                // заставить методы класса использовать 'this'
                exceptMethods: [],
            },
        ],
        complexity: ['off', 4], // Ограничение цикломатической сложности
        'consistent-return': 'error', // return всегда должен что-то отдавать
        curly: ['error', 'multi-line'], // all
        'default-case': 'error', // требует default в switch case
        'default-case-last': 'error', //  Сделать default в операторах switch последними
        'default-param-last': 'off',
        'dot-location': ['error', 'property'],
        'dot-notation': 'off',
        eqeqeq: ['error', 'always'],
        'grouped-accessor-pairs': 'error', // Геттер и сеттер для одного и того же свойства должны быть определены рядом друг с другом
        'guard-for-in': 'off',
        'max-classes-per-file': ['error', 1], // Максимальное количество классов в файле
        'no-alert': 'error',
        'no-caller': 'error', // ????????????????
        'no-case-declarations': 'error', // ????????????????
        'no-constructor-return': 'error', // Запретить возвращаемое значение в конструкторе
        'no-div-regex': 'off', // Запретить регулярные выражения, похожие на деление
        'no-else-return': ['error', { allowElseIf: false }],
        'no-empty-function': [
            'error',
            {
                // ???????????????
                allow: ['arrowFunctions', 'functions', 'methods'],
            },
        ],
        'no-empty-pattern': 'error', // Запретить пустые шаблоны деструктуризации
        'no-eq-null': 'off', // Запретить сравнения с нулевым значением без проверки типов
        'no-eval': 'error', // Запретить использование eval ()
        'no-extend-native': 'error', // Запретить добавление в собственные типы
        'no-extra-bind': 'error', // Запретить ненужную привязку функций
        'no-extra-label': 'error', // Запретить ненужные ярлыки
        'no-fallthrough': 'error', // ?????????????????????
        'no-floating-decimal': 'error',
        'no-global-assign': 'error', // Запретить переназначение глобальных объектов только для чтения
        'no-implicit-coercion': [
            'off',
            {
                // Запретить неявное преобразование типов
                boolean: false,
                number: true,
                string: true,
            },
        ],
        'no-implicit-globals': 'off', // Запретить var и именованные функции в глобальной области видимости
        'no-implied-eval': 'error', // Запретить использование метода eval ()
        'no-invalid-this': 'off',
        'no-iterator': 'error', // Запретить использование свойства __iterator__
        'no-labels': 'error', // ?????????????
        'no-lone-blocks': 'error', // Запретить ненужные вложенные блоки
        'no-loop-func': 'error', // Запретить создание функций внутри циклов
        'no-magic-numbers': 'off',
        'no-multi-spaces': [
            'error',
            {
                // Запретить использование нескольких пробелов
                ignoreEOLComments: true,
            },
        ],
        'no-multi-str': 'error',
        'no-new': 'off', // ???????????
        'no-new-func': 'error',
        'no-new-wrappers': 'error', // Запрещает создание новых экземпляров String, Number и Boolean
        'no-nonoctal-decimal-escape': 'off', // ????????
        'no-octal': 'error', // запретить использование восьмеричных литералов, типа 01, 02 и т.д.
        'no-param-reassign': 'error', // Запретить переназначение параметров функции,
        'no-proto': 'error', // Запретить использование __proto__
        'no-redeclare': 'error', // запретить повторное объявление переменной
        'no-restricted-properties': 'error', // Запретить определенные свойства объекта
        'no-return-assign': ['error', 'always'], // запретить использование присваивания в операторе возврата
        'no-return-await': 'error', //  запретить await в return
        'no-script-url': 'off', // ?????????????
        'no-self-assign': 'error',
        'no-self-compare': 'error', // запретить сравнения, если обе стороны абсолютно одинаковы
        'no-sequences': 'error', // запретить использование оператора запятой
        'no-throw-literal': 'error',
        'no-unmodified-loop-condition': 'off', // ???????????/
        'no-unused-expressions': [
            'error',
            {
                // ????????????
                allowShortCircuit: false,
                allowTernary: true, // Разрешает пример(a ? b : f.c())
                allowTaggedTemplates: false,
            },
        ],
        'no-unused-labels': 'error',
        'no-useless-call': 'off',
        'no-useless-catch': 'error',
        'no-useless-concat': 'error',
        'no-useless-escape': 'error', // ???????????
        'no-useless-return': 'error',
        'no-void': 'error',
        'no-warning-comments': 'off',
        'no-with': 'error',
        'prefer-named-capture-group': 'off',
        'prefer-promise-reject-errors': ['error', { allowEmptyReject: true }], // требовать использования объектов Error в Promise reject()
        'prefer-regex-literals': 'off',
        radix: ['error', 'as-needed'], // требует использования второго аргумента для parseInt ()
        'require-await': 'off',
        'require-unicode-regexp': 'off',
        'vars-on-top': 'error',
        'wrap-iife': 'off',
        yoda: 'error',
        'no-console': 1,
        'no-debugger': 1,
        'no-undef': 0,
    },
};

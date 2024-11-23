module.exports = {
    rules: {
        '@typescript-eslint/consistent-type-imports': 'error', // Требует импорт типов через import type
        '@typescript-eslint/adjacent-overload-signatures': 'warn', //Требовать, чтобы перегрузки элементов были последовательными
        '@typescript-eslint/await-thenable': 'warn', //Запрещает ожидание значения не являющегося .then able
        '@typescript-eslint/ban-tslint-comment': 'error', // запрет комментариев @tslint: <rule>
        '@typescript-eslint/ban-ts-comment': [ //Запрещает использование комментариев @ts-<directive> или требует описания после директивы
            'error',
            {
                'ts-expect-error': 'allow-with-description',
                'ts-ignore': true,
                'ts-nocheck': true,
                'ts-check': false,
                minimumDescriptionLength: 5,
            },
        ],
        '@typescript-eslint/explicit-function-return-type': 'warn', //Требует явно указывать, что возвращает функция
        '@typescript-eslint/ban-types': 'off', //Запрещает использование определенных типов
        '@typescript-eslint/explicit-module-boundary-types': 'off', //Требовать явного возврата и типов аргументов в методах открытого класса экспортируемых функций и классов
        '@typescript-eslint/no-empty-interface': 'error', //Запретить объявление пустых интерфейсов
        '@typescript-eslint/no-explicit-any': 'warn', //Запретить использование any типа
        '@typescript-eslint/no-extra-non-null-assertion': 'error', //Запретить лишнее ненулевое утверждение
        '@typescript-eslint/no-floating-promises': 'off', //Требует, чтобы значения, подобные promise, обрабатывались должным образом
        '@typescript-eslint/no-for-in-array': 'error', //Запретить перебор массива с помощью цикла for-in
        '@typescript-eslint/no-inferrable-types': [ //Запрещает явное объявление типа для переменных или параметров, инициализированных числом, строкой или логическим значением.
            'off',
            {
                ignoreParameters: true,
                ignoreProperties: true,
            },
        ],
        '@typescript-eslint/no-misused-new': 'error', //Обеспечить допустимое определение нового и конструктора
        '@typescript-eslint/no-misused-promises': [//Избегайте использования promise в местах, не предназначенных для их обработки
            'off',
            {
                checksConditionals: true,
            },
        ],
   /*?*/'@typescript-eslint/no-namespace': 'warn', //Запретить использование настраиваемых модулей и пространств имен TypeScript
        '@typescript-eslint/no-non-null-asserted-optional-chain': 'warn', //Запрещает использование ненулевого утверждения после необязательного цепного выражения
        '@typescript-eslint/no-non-null-assertion': 'off', //Запрещает ненулевые утверждения с использованием ! постфиксный оператор
        '@typescript-eslint/no-this-alias': 'off', //Disallow aliasing this
        '@typescript-eslint/no-unnecessary-type-assertion': 'warn', //Предупреждает, если утверждение type не изменяет type выражения
        '@typescript-eslint/no-unsafe-assignment': 'off', //Запрещает any присваивать какие-либо переменным и свойствам
        '@typescript-eslint/no-unsafe-call': 'off', //Запрещает вызывать значение any type
        '@typescript-eslint/no-shadow': 'error', // запретить объявление переменных, уже объявленных во внешней области видимости (баг изза eslint/no-shadow)
        '@typescript-eslint/no-unsafe-member-access': 'off', //Запрещает доступ к членам для any типизированных переменных
        '@typescript-eslint/no-unsafe-return': 'off', //Запрещает возвращать any из функции
        '@typescript-eslint/no-var-requires': 'off', //Запрещает использование операторов require, кроме операторов импорта
        '@typescript-eslint/prefer-as-const': 'error', //Предпочитайте использование as const над буквальным типом
        '@typescript-eslint/prefer-optional-chain': 'warn', // Писать цепные необязательные ?, вместо логических &
        '@typescript-eslint/prefer-namespace-keyword': 'off', //Требовать использования ключевого слова namespace вместо ключевого слова module для объявления пользовательских модулей TypeScript
        '@typescript-eslint/prefer-regexp-exec': 'off', //Обеспечьте использование RegExp # exec вместо String # match, если не указан глобальный флаг
        '@typescript-eslint/promise-function-async': 'warn', // Требует, если функция возвращает промис, то функция async
        '@typescript-eslint/restrict-plus-operands': 'error', // При мат. операциях переменные должны быть одного типа
        '@typescript-eslint/type-annotation-spacing': 'warn', // const t: string - пробел после `:`
        '@typescript-eslint/switch-exhaustiveness-check': 'warn', // Проверяет полный ли перебор делает switch, можно откзаться от полного перебора, указав default
        '@typescript-eslint/restrict-template-expressions': 'off', //Обеспечить, чтобы выражения литералов шаблона имели строковый тип
        '@typescript-eslint/triple-slash-reference': 'warn', //Устанавливает уровень предпочтения для директив с тройной косой чертой по сравнению с объявлениями импорта в стиле ES6
        '@typescript-eslint/unbound-method': 'off', //Принудительно вызывает несвязанные методы с ожидаемой областью видимости
        '@typescript-eslint/unified-signatures': 'warn', // Указывает, когда можно объединить несоклько функций в одну с помощью rest
        //////////////////////////////////////////////////////////////////
        '@typescript-eslint/no-array-constructor': 'off', // Disallow generic Array
        '@typescript-eslint/no-empty-function': 'off', //Это правило расширяет базовое eslint/no-array-constructor правило. Он добавляет поддержку универсально типизированного Array конструктора ( new Array<Foo>()).
        '@typescript-eslint/no-extra-semi': 'warn', //Запретить ненужные точки с запятой
        '@typescript-eslint/no-implied-eval': 'off', //Запретить использование методов, подобных eval ()
        '@typescript-eslint/no-unused-vars': [ //Запретить неиспользуемые переменные
            'warn',
            { varsIgnorePattern: '^_', argsIgnorePattern: '^_' },
        ],
        '@typescript-eslint/require-await': 'warn', //Запретить асинхронные функции, у которых нет выражения ожидания
    },
};

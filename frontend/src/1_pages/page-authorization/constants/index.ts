import type { TFunction } from 'i18next';

import type { ValidationFormType } from 'shared/ui/form/types';

export enum ActionTypes {
    signIn = 'sign-in',
    signUp = 'sign-up',
    recovery = 'recovery',
    code = 'code',
}

export const VALIDATION_SIGN_UP_FORM = (t: TFunction): ValidationFormType[] => [
    {
        name: 'password',
        trigger: 'blur',
        rules: [{
            message: t('pages.auth.valid.password'),
            validator: (val: unknown): boolean => {
                const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@#$%^&+=!-])[A-Za-z\d@#$%^&+=!-]{8,}$/;
                return passwordRegex.test(val?.value);
            },
        }],
    },
];

export const VALIDATION_RECOVERY_FORM = (t: TFunction): ValidationFormType[] => [
    {
        name: 'new-password',
        trigger: 'blur',
        rules: [{
            message: t('pages.auth.valid.password'),
            validator: (val: unknown): boolean => {
                const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@#$%^&+=!-])[A-Za-z\d@#$%^&+=!-]{8,}$/;
                return passwordRegex.test(val?.value);
            },
        }],
    },
];

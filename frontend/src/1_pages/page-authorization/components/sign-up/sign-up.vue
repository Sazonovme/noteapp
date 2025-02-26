<template>
    <div class="header">
        <h2>
            Регистрация
        </h2>
        <RouterLink :to="{ name: 'auth', query: { type: 'signin' } }">
            <ElButton
                size="small"
                plain
            >
                Войти
            </ElButton>
        </RouterLink>
    </div>
    <ElForm
        :ref="state.formRef"
        :model="state.ruleForm"
        :rules="{
            login: [{ required: true, message: 'Это поле необходимо' }],
            password: [
                { required: true, message: 'Это поле необходимо' },
                { min: 7, max: 16, message: 'Длина должна быть от 7 до 16 символов' }
            ]
        }"
        class="form"
        status-icon
    >
        <ElFormItem
            prop="login"
            required
        >
            <ElInput
                v-model="state.ruleForm.login"
                type="text"
                placeholder="Логин"
                size="large"
            />
        </ElFormItem>
        <ElFormItem
            prop="password"
            required
        >
            <ElInput
                v-model="state.ruleForm.password"
                type="password"
                placeholder="Пароль"
                show-password
                size="large"
            />
        </ElFormItem>

        <ElFormItem>
            <ElButton
                size="large"
                type="primary"
                :loading="state.isLoading.value"
                @click="() => actions.onRegistration(state.formRef)"
            >
                Зарегистрироваться
            </ElButton>
        </ElFormItem>
    </ElForm>
</template>

<script setup lang="ts">
import { ElInput, ElForm, ElButton, ElFormItem } from 'element-plus';

import { useSignUp } from './use-sign-up.ts';

const { state, actions } = useSignUp();
</script>

<style scoped src="./sign-up.css" />

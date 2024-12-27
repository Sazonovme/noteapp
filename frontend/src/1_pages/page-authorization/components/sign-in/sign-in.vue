<template>
    <div class="header">
        <h2>
            Авторизация
        </h2>
        <RouterLink :to="{ name: 'auth', query: { type: 'signup' } }">
            <ElButton
                size="small"
                plain
            >
                Регистрация
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
                type="primary"
                size="large"
                :loading="state.isLoading.value"
                @click="actions.onLogin(state.formRef)"
            >
                Войти
            </ElButton>
        </ElFormItem>
    </ElForm>
</template>

<script setup lang="ts">
import { ElInput, ElForm, ElButton, ElFormItem } from 'element-plus';

import { useSignIn } from './use-sign-in.ts';

const { state, actions } = useSignIn();
</script>

<style scoped src="./sign-in.css" />

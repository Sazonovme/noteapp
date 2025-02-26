import { useRouter } from 'vue-router';
import { reactive, ref } from 'vue';
import { FormInstance } from 'element-plus';

import { login } from '@features/user/api';

import { ROUTES_PATH_COMMON } from '@shared/constants';
import { storage, LOCAL_STORAGE_UPDATED_EVENT } from '@shared/utils';

const useState = () => {
    const isLoading = ref(false);
    const loginValue = ref('');
    const passwordValue = ref('');
    const navigate = useRouter();
    const ruleForm = reactive({
        login: '',
        password: '',
    });
    const formRef = ref<FormInstance>();

    return {
        loginValue,
        passwordValue,
        isLoading,
        navigate,
        ruleForm,
        formRef,
    };
};

type UseActionsType = ReturnType<typeof useState>;
const useActions = (state: UseActionsType) => {
    const onLogin = async (form: any): Promise<void> => {
        try {
            if (!form.value) return;
            await form.value.validate(async (valid: any, fields: any) => {
                if (valid) {
                    state.isLoading.value = true;

                    const response = await login({
                        login: form.email.value,
                        password: form.password.value,
                    });

                    storage.local.setItem(LOCAL_STORAGE_UPDATED_EVENT, response.data.token);
                    state.navigate.push(ROUTES_PATH_COMMON.HOME);
                } else {
                    console.log('no valid fields!', fields);
                }
            });
        } catch (e: unknown) {
        // Notification({ message: 'sign-in t()', duration: 3000, type: 'error' });
        } finally {
            state.isLoading.value = false;
        }
    };

    return {
        onLogin,
    };
};

export const useSignIn = () => {
    const state = useState();
    const actions = useActions(state);

    return {
        state,
        actions,
    };
};

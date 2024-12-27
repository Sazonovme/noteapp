import { useRouter } from 'vue-router';
import { reactive, ref } from 'vue';
import { FormInstance } from 'element-plus';

import { registration } from '@features/user/api';

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
    const onRegistration = async (form: any): Promise<void> => {
        try {
            if (!form.value) return;
            await form.value.validate(async (valid: any, fields: any) => {
                if (valid) {
                    state.isLoading.value = true;
                    const response = await registration({
                        login: form.email.value, //TODO потом добавить phone
                        username: form.username.value,
                        password: form.password.value,
                    });

                    storage.local.setItem(LOCAL_STORAGE_UPDATED_EVENT, response.data.token);
                    state.navigate.push(ROUTES_PATH_COMMON.HOME);
                } else {
                    console.log('no valid fields!', fields);
                }
            });
        } catch (e: any) {
        // Notification({ message: e?.data?.message ?? i18next.t('Ошибка при регистрации t()'), duration: 3000, type: 'error' });
        } finally {
            state.isLoading.value = false;
        }
    };

    return {
        onRegistration,
    };
};

export const useSignUp = () => {
    const state = useState();
    const actions = useActions(state);

    return {
        state,
        actions,
    };
};

/*

const rules = reactive<FormRules<RuleForm>>({
  name: [
    { required: true, message: 'Please input Activity name', trigger: 'blur' },
    { min: 3, max: 5, message: 'Length should be 3 to 5', trigger: 'blur' },
  ],
  region: [
    {
      required: true,
      message: 'Please select Activity zone',
      trigger: 'change',
    },
  ],
  count: [
    {
      required: true,
      message: 'Please select Activity count',
      trigger: 'change',
    },
  ],
  date1: [
    {
      type: 'date',
      required: true,
      message: 'Please pick a date',
      trigger: 'change',
    },
  ],
  date2: [
    {
      type: 'date',
      required: true,
      message: 'Please pick a time',
      trigger: 'change',
    },
  ],
  location: [
    {
      required: true,
      message: 'Please select a location',
      trigger: 'change',
    },
  ],
  type: [
    {
      type: 'array',
      required: true,
      message: 'Please select at least one activity type',
      trigger: 'change',
    },
  ],
  resource: [
    {
      required: true,
      message: 'Please select activity resource',
      trigger: 'change',
    },
  ],
  desc: [
    { required: true, message: 'Please input activity form', trigger: 'blur' },
  ],
})
 */

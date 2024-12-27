import { ref, onMounted, onUnmounted } from 'vue';

const useState = () => {
    const backgroundRef = ref<HTMLDivElement>();

    return { backgroundRef };
};

const useActions = (state: ReturnType<typeof useState>) => {
    const onMoveBackground = (event: MouseEvent): void => {
        const height = (innerHeight / 2 - event.pageY) * 0.2;
        const width = (innerWidth / 2 - event.pageX) * 0.2;
        const heightSymbol = height < 0 ? '+' : '-';
        const widthSymbol = width < 0 ? '+' : '-';

        state.backgroundRef.value!.style.top = `calc(-50% ${heightSymbol} ${Math.abs(height)}px)`;
        state.backgroundRef.value!.style.left = `calc(-50% ${widthSymbol} ${Math.abs(width)}px)`;
    };

    return { onMoveBackground };
};

export const usePageAuthorization = () => {
    const state = useState();
    const actions = useActions(state);

    onMounted(() => {
        document.documentElement.addEventListener('mousemove', actions.onMoveBackground);
    });

    onUnmounted(() => {
        document.documentElement.removeEventListener('mousemove', actions.onMoveBackground);
    });

    return {
        ...state,
        ...actions,
    };
};

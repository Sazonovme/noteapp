import { encodeBase64, decodeBase64 } from './base64';

export const LOCAL_STORAGE_UPDATED_EVENT = 'localStorageUpdated' as const;
export const SESSION_STORAGE_UPDATED_EVENT = 'sessionStorageUpdated' as const;
export type EventStorageUpdatedTypes = Event & {
    key: string,
    value: string | null,
};

export const storage = {
    local: {
        ...window.localStorage,
        getItem<T extends string = string>(value: string): T | null { return decodeBase64(window?.localStorage.getItem(encodeBase64(value)) || '') as T; },
        setItem: <T extends string = string>(key: string, value: T): void => {
            const event = new Event(LOCAL_STORAGE_UPDATED_EVENT) as EventStorageUpdatedTypes;
            event.key = key;
            event.value = value;
            window.dispatchEvent(event);
            document.dispatchEvent(event);

            window?.localStorage.setItem(encodeBase64(key), encodeBase64(value));
        },
        removeItem: (key: string): void => {
            const event = new Event(LOCAL_STORAGE_UPDATED_EVENT) as EventStorageUpdatedTypes;
            event.key = key;
            event.value = null;
            window.dispatchEvent(event);
            document.dispatchEvent(event);

            window?.localStorage.removeItem(encodeBase64(key));
        },
    },
    session: {
        ...window.sessionStorage,
        getItem: <T extends string = string>(value: string): T | null => decodeBase64(window?.sessionStorage.getItem(encodeBase64(value)) || '') as T,
        setItem: <T extends string = string>(key: string, value: T): void => {
            const event = new Event(SESSION_STORAGE_UPDATED_EVENT) as EventStorageUpdatedTypes;
            event.key = key;
            event.value = value;
            window.dispatchEvent(event);
            document.dispatchEvent(event);

            window?.sessionStorage.setItem(encodeBase64(key), encodeBase64(value));
        },
        removeItem: (key: string): void => {
            const event = new Event(SESSION_STORAGE_UPDATED_EVENT) as EventStorageUpdatedTypes;
            event.key = key;
            event.value = null;
            window.dispatchEvent(event);
            document.dispatchEvent(event);

            window?.sessionStorage.removeItem(encodeBase64(key));
        },
    },
};

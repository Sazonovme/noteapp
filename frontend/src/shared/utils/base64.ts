export const encodeBase64 = (val: string): string => window.btoa(val);

export const decodeBase64 = (val: string): string => window.atob(val);

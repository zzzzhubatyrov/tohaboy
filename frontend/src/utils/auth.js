// Ключ для хранения токена в localStorage
const TOKEN_KEY = 'token'

// Сохранить токен
export const setToken = (token) => {
    localStorage.setItem(TOKEN_KEY, token);
};

// Получить токен
export const getToken = () => {
    return localStorage.getItem(TOKEN_KEY);
};

// Удалить токен
export const removeToken = () => {
    localStorage.removeItem(TOKEN_KEY);
};

// Проверить наличие токена
export const hasToken = () => {
    return !!getToken();
}; 
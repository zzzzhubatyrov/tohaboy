// Ключ для хранения токена в localStorage
const TOKEN_KEY = 'token'
const USER_KEY = 'user'

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

// Сохранить информацию о пользователе
export const setUser = (user) => {
    localStorage.setItem(USER_KEY, JSON.stringify(user));
};

// Получить информацию о пользователе
export const getUser = () => {
    const userStr = localStorage.getItem(USER_KEY);
    try {
        return userStr ? JSON.parse(userStr) : null;
    } catch {
        return null;
    }
};

// Удалить информацию о пользователе
export const removeUser = () => {
    localStorage.removeItem(USER_KEY);
};

// Очистить все данные аутентификации
export const clearAuth = () => {
    removeToken();
    removeUser();
}; 
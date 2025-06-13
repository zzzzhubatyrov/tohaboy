<script setup>
import {ref} from "vue";
import {useRouter} from "vue-router";
import {Login, Register} from "../../wailsjs/go/service/AuthService.js";
import {setToken} from "../utils/auth";

const router = useRouter();

const form = ref({
  username: '',
  password: ''
});

const isLogin = ref(true);
const showPassword = ref(false);
const loading = ref(false);
const errors = ref({
  username: '',
  password: ''
});

const notification = ref({
  show: false,
  message: '',
  type: 'success'
});

async function handleSubmit(e) {
  e.preventDefault();
  clearErrors();
  
  if (!validateForm()) {
    return;
  }

  loading.value = true;
  console.log('Attempting auth with:', { ...form.value, password: '[REDACTED]' });

  try {
    if (isLogin.value) {
      console.log('Logging in...');
      const response = await Login(form.value);
      console.log('Login response:', response);
      
      if (response && response.token) {
        setToken(response.token);
        showNotification('Успешный вход', 'success');
        router.push('/');
      } else {
        throw new Error('Неверный логин или пароль');
      }
    } else {
      console.log('Registering...');
      const response = await Register(form.value);
      console.log('Register response:', response);
      
      if (response && response.id) {
        showNotification('Регистрация успешна', 'success');
        isLogin.value = true;
        form.value.password = '';
      } else {
        throw new Error('Ошибка при регистрации');
      }
    }
  } catch (error) {
    console.error('Auth error:', error);
    showNotification(
      error.message || 'Произошла ошибка при авторизации',
      'error'
    );
  } finally {
    loading.value = false;
  }
}

function showNotification(message, type = 'success') {
  notification.value = {
    show: true,
    message,
    type
  };
  setTimeout(() => {
    notification.value.show = false;
  }, 4000);
}

function toggleMode() {
  isLogin.value = !isLogin.value;
  form.value = { username: '', password: '' };
  clearErrors();
}

function clearErrors() {
  errors.value = {
    username: '',
    password: ''
  };
}

function validateForm() {
  let isValid = true;
  errors.value = {
    username: '',
    password: ''
  };

  if (!form.value.username.trim()) {
    errors.value.username = 'Введите имя пользователя';
    isValid = false;
  }

  if (!form.value.password) {
    errors.value.password = 'Введите пароль';
    isValid = false;
  } else if (form.value.password.length < 4) {
    errors.value.password = 'Пароль должен быть не менее 4 символов';
    isValid = false;
  }

  return isValid;
}
</script>

<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <h1 class="auth-title">{{ isLogin ? 'Вход' : 'Регистрация' }}</h1>
        <p class="auth-subtitle">
          {{ isLogin ? 'Войдите в систему для продолжения' : 'Создайте новую учетную запись' }}
        </p>
      </div>

      <form @submit="handleSubmit" class="auth-form">
        <div class="form-group">
          <label for="username">Имя пользователя</label>
          <input
              id="username"
              v-model="form.username"
              type="text"
              required
              class="form-input"
              :class="{ 'error': errors.username }"
              placeholder="Введите имя пользователя"
          />
          <span v-if="errors.username" class="error-text">{{ errors.username }}</span>
        </div>

        <div class="form-group">
          <label for="password">Пароль</label>
          <div class="password-input">
            <input
                id="password"
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                required
                class="form-input"
                :class="{ 'error': errors.password }"
                placeholder="Введите пароль"
            />
            <button
                type="button"
                class="password-toggle"
                @click="showPassword = !showPassword"
            >
              <svg v-if="showPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                <line x1="1" y1="1" x2="23" y2="23"/>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
            </button>
          </div>
          <span v-if="errors.password" class="error-text">{{ errors.password }}</span>
        </div>

        <button type="submit" class="btn btn-primary" :disabled="loading">
          <span v-if="loading" class="spinner"></span>
          {{ isLogin ? 'Войти' : 'Зарегистрироваться' }}
        </button>

        <div class="auth-toggle">
          {{ isLogin ? 'Нет аккаунта?' : 'Уже есть аккаунт?' }}
          <button type="button" class="btn-link" @click="toggleMode">
            {{ isLogin ? 'Зарегистрироваться' : 'Войти' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Notification -->
    <div v-if="notification.show" :class="['notification', 'notification-' + notification.type]">
      {{ notification.message }}
    </div>
  </div>
</template>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: #f8fafc;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  padding: 32px;
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.auth-title {
  font-size: 24px;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 8px;
}

.auth-subtitle {
  font-size: 14px;
  color: #64748b;
  margin: 0;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #475569;
}

.form-input {
  padding: 10px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  color: #1e293b;
  background: white;
  width: 100%;
  transition: all 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-input.error {
  border-color: #ef4444;
}

.error-text {
  font-size: 12px;
  color: #ef4444;
}

.password-input {
  position: relative;
}

.password-toggle {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  padding: 0;
  color: #94a3b8;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.password-toggle:hover {
  color: #64748b;
}

.password-toggle svg {
  width: 20px;
  height: 20px;
  stroke-width: 2;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 24px;
  border-radius: 6px;
  font-weight: 500;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn-link {
  background: none;
  color: #3b82f6;
  padding: 0;
  font-weight: 500;
}

.btn-link:hover {
  color: #2563eb;
  text-decoration: underline;
}

.auth-toggle {
  text-align: center;
  font-size: 14px;
  color: #64748b;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid #ffffff;
  border-top: 2px solid transparent;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.notification {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 16px 24px;
  border-radius: 8px;
  color: white;
  font-weight: 500;
  box-shadow: 0 10px 25px -3px rgba(0, 0, 0, 0.1);
  z-index: 1001;
  animation: slideIn 0.3s ease-out;
}

.notification-success {
  background: #10b981;
}

.notification-error {
  background: #ef4444;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@media (max-width: 480px) {
  .auth-card {
    padding: 24px;
  }

  .auth-title {
    font-size: 20px;
  }

  .btn {
    padding: 10px 20px;
  }
}
</style>
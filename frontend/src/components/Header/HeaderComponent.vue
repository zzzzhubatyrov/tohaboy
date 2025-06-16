<script setup>
import { useRoute, useRouter } from 'vue-router'
import { ref, onMounted } from 'vue'
import { clearAuth, getToken, getUser } from '../../utils/auth'

const route = useRoute()
const router = useRouter()

const navLeftLinks = [
  {to: '/', label: 'Оборудование', activeFor: '/'},
  {to: '/suppliers', label: 'Поставщики', activeFor: '/suppliers'},
  {to: '/documents', label: 'Документы', activeFor: '/documents'},
  {to: '/movements', label: 'Передача', activeFor: '/movements'},
]

const username = ref('')

onMounted(() => {
  // Проверяем наличие токена
  if (!getToken()) {
    router.push('/auth')
    return
  }

  // Получаем информацию о пользователе из хранилища
  const user = getUser()
  if (user) {
    username.value = user.username
  } else {
    clearAuth()
    router.push('/auth')
  }
})

function logout() {
  clearAuth()
  router.push('/auth')
}
</script>

<template>
  <header class="header">
    <nav class="nav">
      <router-link
          v-for="link in navLeftLinks"
          :key="link.to"
          :to="link.to"
          class="nav-link"
          :class="{ 'active': route.path === link.activeFor }"
      >
        {{ link.label }}
      </router-link>
    </nav>

    <div class="user-menu">
      <button class="user-button" @click="logout">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
          <polyline points="16 17 21 12 16 7"/>
          <line x1="21" y1="12" x2="9" y2="12"/>
        </svg>
        <span class="username">{{ username || 'Выйти' }}</span>
      </button>
    </div>
  </header>
</template>

<style scoped>
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: white;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  z-index: 1000;
}

.nav {
  display: flex;
  gap: 8px;
  height: 100%;
  align-items: center;
}

.nav-link {
  text-decoration: none;
  color: #475569;
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 500;
  transition: all 0.2s;
  height: 40px;
  display: flex;
  align-items: center;
}

.nav-link:hover {
  background: #f1f5f9;
  color: #1e293b;
}

.nav-link.active {
  background: #e0f2fe;
  color: #0369a1;
}

.user-menu {
  display: flex;
  align-items: center;
}

.user-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border: none;
  background: none;
  color: #475569;
  font-weight: 500;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s;
}

.user-button:hover {
  background: #f1f5f9;
  color: #1e293b;
}

.icon {
  width: 16px;
  height: 16px;
  stroke-width: 2;
}

.username {
  font-size: 14px;
}

@media (max-width: 768px) {
  .header {
    padding: 0 16px;
  }

  .nav-link {
    padding: 8px 12px;
    font-size: 14px;
  }

  .user-button {
    padding: 8px 12px;
  }

  .username {
    display: none;
  }
}
</style>
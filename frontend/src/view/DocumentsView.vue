<script setup>

import HeaderComponent from "../components/Header/HeaderComponent.vue";
import {onMounted, ref, computed} from "vue";
import {
  CreateDocument,
  DeleteDocument,
  GetAllDocuments,
  UpdateDocument,
  ApproveDocument
} from "../../wailsjs/go/service/DocumentService";
import {
  GetAllLocations,
} from "../../wailsjs/go/service/LocationService";
import {
  GetAllEquipment,
} from "../../wailsjs/go/service/EquipmentService";
import {
  Login,
} from "../../wailsjs/go/service/AuthService";
import {
  GetUser,
} from "../../wailsjs/go/service/UserService";

// Состояние
const documents = ref([])
const locations = ref([])
const equipmentList = ref([])
const loading = ref(false)
const showModal = ref(false)
const modalMode = ref('create') // 'create', 'edit', 'view'
const currentUser = ref(null)
const currentDocument = ref(null)
const selectedDocument = ref(null)

// Фильтры и поиск
const searchQuery = ref('')
const typeFilter = ref('')
const statusFilter = ref('')

// Сортировка
const sortField = ref('date')
const sortDirection = ref('desc')

// Пагинация
const currentPage = ref(1)
const itemsPerPage = ref(20)

// Уведомления
const notification = ref({
  show: false,
  message: '',
  type: 'success'
})

// Вычисляемые свойства
const filteredDocuments = computed(() => {
  let filtered = [...documents.value]

  // Поиск
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(doc =>
        doc.number.toLowerCase().includes(query)
    )
  }

  // Фильтр по типу
  if (typeFilter.value) {
    filtered = filtered.filter(doc => doc.type === typeFilter.value)
  }

  // Фильтр по статусу
  if (statusFilter.value) {
    filtered = filtered.filter(doc => doc.status === statusFilter.value)
  }

  return filtered
})

const sortedDocuments = computed(() => {
  return [...filteredDocuments.value].sort((a, b) => {
    let aVal = a[sortField.value]
    let bVal = b[sortField.value]

    if (typeof aVal === 'string') {
      aVal = aVal.toLowerCase()
      bVal = bVal.toLowerCase()
    }

    if (sortField.value === 'date') {
      aVal = new Date(aVal)
      bVal = new Date(bVal)
    }

    if (sortDirection.value === 'asc') {
      return aVal < bVal ? -1 : aVal > bVal ? 1 : 0
    } else {
      return aVal > bVal ? -1 : aVal < bVal ? 1 : 0
    }
  })
})

const paginatedDocuments = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return sortedDocuments.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.ceil(sortedDocuments.value.length / itemsPerPage.value)
})

const totalAmount = computed(() => {
  return currentDocument.value.items.reduce((sum, item) => sum + (item.total_price || 0), 0)
})

const totalDocuments = computed(() => documents.value.length)

const draftDocuments = computed(() => {
  return documents.value.filter(doc => doc.status === 'draft').length
})

// Методы
async function loadData() {
  loading.value = true
  try {
    await Promise.all([
      loadDocuments(),
      loadLocations(),
      loadEquipment(),
      loadCurrentUser()
    ])
  } finally {
    loading.value = false
  }
}

async function loadDocuments() {
  try {
    const response = await GetAllDocuments()
    if (response.model) {
      documents.value = response.model
    } else {
      showNotification(response.msg || 'Ошибка загрузки документов', 'error')
    }
  } catch (error) {
    showNotification('Ошибка подключения к серверу', 'error')
  }
}

async function loadLocations() {
  try {
    const response = await GetAllLocations()
    if (response.model) {
      locations.value = response.model
    }
  } catch (error) {
    console.error('Ошибка загрузки локаций:', error)
  }
}

async function loadEquipment() {
  try {
    const response = await GetAllEquipment()
    if (response.model) {
      equipmentList.value = response.model
    }
  } catch (error) {
    console.error('Ошибка загрузки оборудования:', error)
  }
}

async function loadCurrentUser() {
  try {
    // Используем GetUser для получения пользователя admin
    const response = await GetUser("admin")
    if (response.model) {
      currentUser.value = response.model
    } else {
      // Если не удалось получить admin, пробуем залогиниться
      const loginResponse = await Login({ username: "admin", password: "admin" })
      if (loginResponse) {
        currentUser.value = loginResponse
      } else {
        throw new Error('Не удалось получить текущего пользователя')
      }
    }
  } catch (error) {
    console.error('Ошибка загрузки текущего пользователя:', error)
    showNotification('Ошибка загрузки текущего пользователя', 'error')
  }
}

function getEmptyDocument() {
  return {
    id: 0,
    type: 'inventory',
    number: '',
    date: new Date().toISOString().split('T')[0],
    status: 'draft',
    comment: '',
    location_id: '',
    created_by_id: currentUser.value?.id || 0,
    items: []
  }
}

function getEmptyDocumentItem() {
  return {
    equipment_id: '',
    quantity: 1,
    actual_quantity: 0,
    price: 0,
    total_price: 0
  }
}

function generateDocumentNumber() {
  const date = new Date()
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const typePrefix = {
    'inventory': 'ИНВ',
    'transfer': 'ПЕР',
    'write_off': 'СПС',
    'acceptance': 'ПРМ'
  }[currentDocument.value.type]
  
  currentDocument.value.number = `${typePrefix}-${year}${month}-${Math.floor(Math.random() * 1000).toString().padStart(3, '0')}`
}

function updateItemPrice(item) {
  if (item.equipment_id) {
    const equipment = equipmentList.value.find(eq => eq.id === item.equipment_id)
    if (equipment) {
      item.price = equipment.price
      updateTotalPrice(item)
    }
  }
}

function updateTotalPrice(item) {
  item.total_price = item.quantity * item.price
}

function addDocumentItem() {
  currentDocument.value.items.push(getEmptyDocumentItem())
}

function removeDocumentItem(index) {
  currentDocument.value.items.splice(index, 1)
}

function openCreateModal() {
  modalMode.value = 'create'
  currentDocument.value = getEmptyDocument()
  generateDocumentNumber()
  showModal.value = true
}

function viewDocument(document) {
  if (!document) return
  modalMode.value = 'view'
  currentDocument.value = { ...document }
  showModal.value = true
}

function editDocument(document) {
  if (!document) return
  modalMode.value = 'edit'
  currentDocument.value = { ...document }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  currentDocument.value = getEmptyDocument()
}

async function saveDocument() {
  try {
    // Добавляем ID создателя документа
    if (modalMode.value === 'create') {
      currentDocument.value.created_by_id = currentUser.value?.id
    }

    // Проверяем наличие позиций
    if (!currentDocument.value.items.length) {
      showNotification('Добавьте хотя бы одну позицию в документ', 'error')
      return
    }

    // Проверяем обязательные поля
    if (!currentDocument.value.location_id) {
      showNotification('Выберите местоположение', 'error')
      return
    }

    // Проверяем позиции документа
    for (const item of currentDocument.value.items) {
      if (!item.equipment_id) {
        showNotification('Выберите оборудование для всех позиций', 'error')
        return
      }
      if (item.quantity <= 0) {
        showNotification('Количество должно быть больше нуля', 'error')
        return
      }
      if (currentDocument.value.type === 'inventory' && item.actual_quantity < 0) {
        showNotification('Фактическое количество не может быть отрицательным', 'error')
        return
      }
    }

    let response
    if (modalMode.value === 'create') {
      response = await CreateDocument(currentDocument.value)
    } else {
      response = await UpdateDocument(currentDocument.value)
    }

    if (response.model) {
      showNotification(response.msg || 'Документ успешно сохранен', 'success')
      closeModal()
      await loadDocuments()
    } else {
      showNotification(response.msg || 'Ошибка сохранения', 'error')
    }
  } catch (error) {
    showNotification('Ошибка подключения к серверу', 'error')
  }
}

async function deleteDocument(document) {
  if (!document) return
  if (!confirm(`Удалить документ "${document.number}"?`)) return

  try {
    const response = await DeleteDocument(document.id)
    if (response.model) {
      showNotification(response.msg, 'success')
      await loadDocuments()
    } else {
      showNotification(response.msg || 'Ошибка удаления', 'error')
    }
  } catch (error) {
    showNotification('Ошибка подключения к серверу', 'error')
  }
}

async function approveDocument(document) {
  if (!document) return
  if (!confirm(`Утвердить документ "${document.number}"?`)) return

  try {
    const response = await ApproveDocument(document.id, currentUser.value?.id)
    if (response.model) {
      showNotification(response.msg || 'Документ успешно утвержден', 'success')
      await loadDocuments()
    } else {
      showNotification(response.msg || 'Ошибка утверждения', 'error')
    }
  } catch (error) {
    showNotification('Ошибка подключения к серверу', 'error')
  }
}

function showNotification(message, type = 'success') {
  notification.value = {
    show: true,
    message,
    type
  }

  setTimeout(() => {
    notification.value.show = false
  }, 4000)
}

function selectDocument(document) {
  selectedDocument.value = document
}

function sortBy(field) {
  if (sortField.value === field) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortField.value = field
    sortDirection.value = 'asc'
  }
}

function getDocumentTypeText(type) {
  const typeMap = {
    'inventory': 'Акт описи',
    'transfer': 'Акт перемещения',
    'write_off': 'Акт списания',
    'acceptance': 'Акт приемки'
  }
  return typeMap[type] || type
}

function getStatusText(status) {
  const statusMap = {
    'draft': 'Черновик',
    'completed': 'Утвержден',
    'canceled': 'Отменен'
  }
  return statusMap[status] || status
}

function formatDate(date) {
  return new Date(date).toLocaleDateString('ru-RU')
}

function formatPrice(price) {
  if (!price) return '—'
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB'
  }).format(price)
}

function applyFilters() {
  currentPage.value = 1
}

function refreshData() {
  loadData()
}

// Инициализация после загрузки компонента
onMounted(async () => {
  await loadData()
  currentDocument.value = getEmptyDocument()
})
</script>

<template>
  <HeaderComponent />
  <div class="document-container">
    <!-- Header -->
    <div class="header">
      <h1 class="title">Управление документами</h1>
      <div class="header-stats">
        <div class="stat-item">
          <span class="stat-number">{{ totalDocuments }}</span>
          <span class="stat-label">Всего документов</span>
        </div>
        <div class="stat-item">
          <span class="stat-number">{{ draftDocuments }}</span>
          <span class="stat-label">Черновики</span>
        </div>
      </div>
    </div>

    <!-- Action Panel -->
    <div class="action-panel">
      <div class="action-left">
        <div class="action-buttons-group">
          <button class="btn btn-primary" @click="openCreateModal">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M12 5v14M5 12h14"/>
            </svg>
            Создать документ
          </button>

          <button class="btn btn-secondary" @click="viewDocument(selectedDocument)" :disabled="!selectedDocument">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
              <circle cx="12" cy="12" r="3"/>
            </svg>
            Просмотр
          </button>

          <button class="btn btn-secondary" @click="editDocument(selectedDocument)" 
                  :disabled="!selectedDocument || selectedDocument.status !== 'draft'">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="m18.5 2.5 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
            Редактировать
          </button>

          <button class="btn btn-success" @click="approveDocument(selectedDocument)"
                  :disabled="!selectedDocument || selectedDocument.status !== 'draft'">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M20 6L9 17l-5-5"/>
            </svg>
            Утвердить
          </button>

          <button class="btn btn-danger" @click="deleteDocument(selectedDocument)" 
                  :disabled="!selectedDocument || selectedDocument.status !== 'draft'">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline points="3,6 5,6 21,6"/>
              <path d="m19,6v14a2,2 0 0,1 -2,2H7a2,2 0 0,1 -2,-2V6m3,0V4a2,2 0 0,1 2,-2h4a2,2 0 0,1 2,2v2"/>
            </svg>
            Удалить
          </button>
        </div>

        <div class="filter-group">
          <select v-model="typeFilter" @change="applyFilters" class="select">
            <option value="">Все типы</option>
            <option value="inventory">Акт описи</option>
            <option value="transfer">Акт перемещения</option>
            <option value="write_off">Акт списания</option>
            <option value="acceptance">Акт приемки</option>
          </select>

          <select v-model="statusFilter" @change="applyFilters" class="select">
            <option value="">Все статусы</option>
            <option value="draft">Черновик</option>
            <option value="completed">Утвержден</option>
            <option value="canceled">Отменен</option>
          </select>
        </div>
      </div>

      <div class="action-right">
        <div class="search-box">
          <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <circle cx="11" cy="11" r="8"/>
            <path d="m21 21-4.35-4.35"/>
          </svg>
          <input
              v-model="searchQuery"
              @input="applyFilters"
              placeholder="Поиск по номеру документа..."
              class="search-input"
          />
        </div>

        <button class="btn btn-secondary" @click="refreshData">
          <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/>
            <path d="M21 3v5h-5"/>
            <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/>
            <path d="M3 21v-5h5"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Documents Table -->
    <div class="table-container">
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <span>Загрузка данных...</span>
      </div>

      <table v-else class="document-table">
        <thead>
        <tr>
          <th @click="sortBy('number')" class="sortable">
            Номер
            <svg v-if="sortField === 'number'" class="sort-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline :points="sortDirection === 'asc' ? '6 9 12 15 18 9' : '6 15 12 9 18 15'"/>
            </svg>
          </th>
          <th>Тип</th>
          <th @click="sortBy('date')" class="sortable">
            Дата
            <svg v-if="sortField === 'date'" class="sort-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline :points="sortDirection === 'asc' ? '6 9 12 15 18 9' : '6 15 12 9 18 15'"/>
            </svg>
          </th>
          <th>Статус</th>
          <th>Местоположение</th>
          <th>Создал</th>
          <th>Утвердил</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="document in paginatedDocuments"
            :key="document.id"
            class="table-row"
            :class="{ 'selected': selectedDocument && selectedDocument.id === document.id }"
            @click="selectDocument(document)">
          <td class="document-number">{{ document.number }}</td>
          <td>{{ getDocumentTypeText(document.type) }}</td>
          <td>{{ formatDate(document.date) }}</td>
          <td>
            <span :class="['status-badge', `status-${document.status}`]">
              {{ getStatusText(document.status) }}
            </span>
          </td>
          <td>{{ document.location?.name || '—' }}</td>
          <td>{{ document.created_by?.username || '—' }}</td>
          <td>{{ document.approved_by?.username || '—' }}</td>
        </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination">
        <button
            @click="currentPage = 1"
            :disabled="currentPage === 1"
            class="pagination-btn"
        >
          Первая
        </button>
        <button
            @click="currentPage--"
            :disabled="currentPage === 1"
            class="pagination-btn"
        >
          ←
        </button>

        <span class="pagination-info">
          {{ currentPage }} из {{ totalPages }}
        </span>

        <button
            @click="currentPage++"
            :disabled="currentPage === totalPages"
            class="pagination-btn"
        >
          →
        </button>
        <button
            @click="currentPage = totalPages"
            :disabled="currentPage === totalPages"
            class="pagination-btn"
        >
          Последняя
        </button>
      </div>
    </div>

    <!-- Document Modal -->
    <div v-if="showModal" class="modal-overlay" @click="closeModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h2>{{ modalMode === 'create' ? 'Создать документ' : modalMode === 'edit' ? 'Редактировать документ' : 'Просмотр документа' }}</h2>
          <button @click="closeModal" class="modal-close">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <form @submit.prevent="saveDocument">
            <div class="form-grid">
              <div class="form-group">
                <label>Тип документа *</label>
                <select
                    v-model="currentDocument.type"
                    :disabled="modalMode === 'view'"
                    required
                    class="form-select"
                    @change="generateDocumentNumber"
                >
                  <option value="inventory">Акт описи</option>
                  <option value="transfer">Акт перемещения</option>
                  <option value="write_off">Акт списания</option>
                  <option value="acceptance">Акт приемки</option>
                </select>
              </div>

              <div class="form-group">
                <label>Номер документа</label>
                <div v-if="modalMode === 'view'" class="form-static-value">
                  {{ currentDocument.number }}
                </div>
                <input
                    v-else
                    v-model="currentDocument.number"
                    disabled
                    class="form-input"
                />
              </div>

              <div class="form-group">
                <label>Дата *</label>
                <input
                    v-model="currentDocument.date"
                    :disabled="modalMode === 'view'"
                    type="date"
                    required
                    class="form-input"
                />
              </div>

              <div class="form-group">
                <label>Местоположение *</label>
                <div v-if="modalMode === 'view'" class="form-static-value">
                  {{ currentDocument.location?.name || '—' }}
                </div>
                <select
                    v-else
                    v-model="currentDocument.location_id"
                    required
                    class="form-select"
                >
                  <option value="">Выберите местоположение</option>
                  <option v-for="location in locations" :key="location.id" :value="location.id">
                    {{ location.name }}
                  </option>
                </select>
              </div>
            </div>

            <div class="form-group">
              <label>Комментарий</label>
              <textarea
                  v-model="currentDocument.comment"
                  :disabled="modalMode === 'view'"
                  rows="2"
                  class="form-textarea"
              ></textarea>
            </div>

            <!-- Document Items -->
            <div class="document-items">
              <h3>Позиции документа</h3>
              
              <div v-if="modalMode !== 'view'" class="add-item-button">
                <button type="button" @click="addDocumentItem" class="btn btn-secondary">
                  <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                    <path d="M12 5v14M5 12h14"/>
                  </svg>
                  Добавить позицию
                </button>
              </div>

              <table class="items-table">
                <thead>
                <tr>
                  <th>Оборудование</th>
                  <th>Количество</th>
                  <th v-if="currentDocument.type === 'inventory'">Факт</th>
                  <th>Цена</th>
                  <th>Сумма</th>
                  <th v-if="modalMode !== 'view'">Действия</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(item, index) in currentDocument.items" :key="index">
                  <td>
                    <div v-if="modalMode === 'view'" class="form-static-value">
                      {{ item.equipment?.name || '—' }}
                    </div>
                    <select
                        v-else
                        v-model="item.equipment_id"
                        class="form-select"
                        @change="updateItemPrice(item)"
                    >
                      <option value="">Выберите оборудование</option>
                      <option v-for="equipment in equipmentList" :key="equipment.id" :value="equipment.id">
                        {{ equipment.name }} ({{ equipment.serial_number }})
                      </option>
                    </select>
                  </td>
                  <td>
                    <div v-if="modalMode === 'view'" class="form-static-value">
                      {{ item.quantity }}
                    </div>
                    <input
                        v-else
                        v-model.number="item.quantity"
                        type="number"
                        min="1"
                        class="form-input"
                        @input="updateTotalPrice(item)"
                    />
                  </td>
                  <td v-if="currentDocument.type === 'inventory'">
                    <div v-if="modalMode === 'view'" class="form-static-value">
                      {{ item.actual_quantity }}
                    </div>
                    <input
                        v-else
                        v-model.number="item.actual_quantity"
                        type="number"
                        min="0"
                        class="form-input"
                    />
                  </td>
                  <td>
                    <div v-if="modalMode === 'view'" class="form-static-value">
                      {{ formatPrice(item.price) }}
                    </div>
                    <input
                        v-else
                        v-model.number="item.price"
                        type="number"
                        min="0"
                        step="0.01"
                        class="form-input"
                        @input="updateTotalPrice(item)"
                    />
                  </td>
                  <td>{{ formatPrice(item.total_price) }}</td>
                  <td v-if="modalMode !== 'view'">
                    <button type="button" @click="removeDocumentItem(index)" class="btn btn-icon delete-btn">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                        <path d="M18 6L6 18M6 6l12 12"/>
                      </svg>
                    </button>
                  </td>
                </tr>
                </tbody>
                <tfoot>
                <tr>
                  <td colspan="3">Итого:</td>
                  <td colspan="2">{{ formatPrice(totalAmount) }}</td>
                  <td v-if="modalMode !== 'view'"></td>
                </tr>
                </tfoot>
              </table>
            </div>

            <div v-if="modalMode !== 'view'" class="modal-actions">
              <button type="button" @click="closeModal" class="btn btn-secondary">
                Отмена
              </button>
              <button type="submit" class="btn btn-primary">
                {{ modalMode === 'create' ? 'Создать' : 'Сохранить' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Notification -->
    <div v-if="notification.show" :class="['notification', `notification-${notification.type}`]">
      {{ notification.message }}
    </div>
  </div>
</template>

<style scoped>
.document-container {
  padding: 24px;
  background: #f8fafc;
  min-height: 100vh;
}

/* Header */
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.title {
  font-size: 32px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.header-stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 24px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.stat-number {
  font-size: 24px;
  font-weight: 700;
  color: #3b82f6;
}

.stat-label {
  font-size: 14px;
  color: #64748b;
  margin-top: 4px;
}

/* Action Panel */
.action-panel {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.action-left {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.action-buttons-group {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  gap: 12px;
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border-radius: 8px;
  border: none;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn-secondary {
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.btn-secondary:hover:not(:disabled) {
  background: #e2e8f0;
}

.btn-success {
  background: #10b981;
  color: white;
}

.btn-success:hover:not(:disabled) {
  background: #059669;
}

.btn-danger {
  background: #ef4444;
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background: #dc2626;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  border: none;
  background: #f8fafc;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-icon:hover {
  background: #e2e8f0;
  color: #475569;
}

.btn-icon.delete-btn:hover {
  background: #fef2f2;
  color: #dc2626;
}

.icon {
  width: 16px;
  height: 16px;
  stroke-width: 2;
}

/* Form Elements */
.select, .search-input {
  padding: 8px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: white;
  font-size: 14px;
  color: #334155;
}

.select:focus, .search-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  width: 16px;
  height: 16px;
  color: #94a3b8;
  stroke-width: 2;
}

.search-input {
  padding-left: 40px;
  width: 280px;
}

/* Table */
.table-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.document-table {
  width: 100%;
  border-collapse: collapse;
}

.document-table th {
  background: #f8fafc;
  padding: 16px 12px;
  text-align: left;
  font-weight: 600;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
  font-size: 14px;
}

.document-table th.sortable {
  cursor: pointer;
  user-select: none;
  position: relative;
}

.document-table th.sortable:hover {
  background: #f1f5f9;
}

.sort-icon {
  width: 14px;
  height: 14px;
  margin-left: 4px;
  vertical-align: middle;
}

.document-table td {
  padding: 16px 12px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  font-size: 14px;
}

.table-row {
  cursor: pointer;
  transition: background-color 0.2s;
}

.table-row:hover {
  background: #f8fafc;
}

.table-row.selected {
  background: #e0f2fe;
}

.table-row.selected:hover {
  background: #bae6fd;
}

.document-number {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: #475569;
}

.status-badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-draft {
  background: #fef3c7;
  color: #92400e;
}

.status-completed {
  background: #dcfce7;
  color: #166534;
}

.status-canceled {
  background: #fee2e2;
  color: #991b1b;
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 1000px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px;
  border-bottom: 1px solid #f1f5f9;
}

.modal-header h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: #1e293b;
}

.modal-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: #f8fafc;
  border-radius: 6px;
  cursor: pointer;
  color: #64748b;
  transition: all 0.2s;
}

.modal-close:hover {
  background: #e2e8f0;
  color: #475569;
}

.modal-body {
  padding: 24px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-weight: 600;
  color: #374151;
  font-size: 14px;
}

.form-input, .form-select, .form-textarea {
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  color: #374151;
  background: white;
  transition: border-color 0.2s;
}

.form-input:focus, .form-select:focus, .form-textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-input:disabled, .form-select:disabled, .form-textarea:disabled {
  background: #f9fafb;
  color: #6b7280;
  cursor: not-allowed;
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.form-static-value {
  padding: 10px 12px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  color: #475569;
  min-height: 40px;
  display: flex;
  align-items: center;
}

/* Document Items */
.document-items {
  margin-top: 24px;
  border-top: 1px solid #f1f5f9;
  padding-top: 24px;
}

.document-items h3 {
  margin: 0 0 16px;
  font-size: 18px;
  color: #1e293b;
}

.add-item-button {
  margin-bottom: 16px;
}

.items-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 16px;
}

.items-table th {
  background: #f8fafc;
  padding: 12px;
  text-align: left;
  font-weight: 600;
  color: #475569;
  font-size: 14px;
  border-bottom: 1px solid #e2e8f0;
}

.items-table td {
  padding: 12px;
  border-bottom: 1px solid #f1f5f9;
}

.items-table tfoot td {
  font-weight: 600;
  color: #1e293b;
  background: #f8fafc;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #f1f5f9;
}

/* Notification */
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

/* Responsive */
@media (max-width: 1200px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .document-container {
    padding: 16px;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .action-panel {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .action-left, .action-right {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .search-input {
    width: 100%;
  }

  .document-table {
    font-size: 12px;
  }

  .document-table th,
  .document-table td {
    padding: 8px 6px;
  }

  .modal {
    width: 95%;
    margin: 10px;
  }

  .modal-header, .modal-body {
    padding: 16px;
  }
}
</style>
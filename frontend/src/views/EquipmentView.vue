<template>
  <HeaderComponent />
  <div class="equipment-container">
    <!-- Header -->
    <div class="header">
      <h1 class="title">Инвентаризация и управление оборудованием</h1>
      <div class="header-stats">
        <div class="stat-item">
          <span class="stat-number">{{ totalItems }}</span>
          <span class="stat-label">Всего позиций</span>
        </div>
        <div class="stat-item">
          <span class="stat-number">{{ formatPrice(totalSum) }}</span>
          <span class="stat-label">Общая стоимость</span>
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
            Добавить оборудование
          </button>

          <button class="btn btn-secondary" @click="exportToExcel">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M14 3v4a1 1 0 0 0 1 1h4"/>
              <path d="M17 21H7a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h7l5 5v11a2 2 0 0 1-2 2z"/>
              <path d="m9 17 2-2m0 0 2-2m-2 2-2-2m2 2 2 2"/>
            </svg>
            Экспорт в Excel
          </button>

          <button class="btn btn-secondary" @click="viewEquipment(selectedEquipment)" :disabled="!selectedEquipment">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
              <circle cx="12" cy="12" r="3"/>
            </svg>
            Просмотр
          </button>

          <button class="btn btn-secondary" @click="editEquipment(selectedEquipment)" :disabled="!selectedEquipment">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="m18.5 2.5 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
            Редактировать
          </button>

          <button class="btn btn-danger" @click="deleteEquipment(selectedEquipment)" :disabled="!selectedEquipment">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline points="3,6 5,6 21,6"/>
              <path d="m19,6v14a2,2 0 0,1 -2,2H7a2,2 0 0,1 -2,-2V6m3,0V4a2,2 0 0,1 2,-2h4a2,2 0 0,1 2,2v2"/>
            </svg>
            Удалить
          </button>

          <button class="btn btn-info" @click="openTransferModal(selectedEquipment)" :disabled="!selectedEquipment">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M16 3h5v5"/>
              <path d="M4 20L21 3"/>
              <path d="M21 16v5h-5"/>
              <path d="M4 4l5 5"/>
            </svg>
            Передать
          </button>
        </div>

        <div class="filter-group">
          <select v-model="categoryFilter" @change="applyFilters" class="select">
            <option value="">Все категории</option>
            <option v-for="category in categories" :key="category.id" :value="category.id">
              {{ category.name }}
            </option>
          </select>

          <select v-model="statusFilter" @change="applyFilters" class="select">
            <option value="">Все статусы</option>
            <option value="available">Доступно</option>
            <option value="in_use">Используется</option>
            <option value="maintenance">На обслуживании</option>
            <option value="written_off">Списано</option>
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
              placeholder="Поиск по названию или серийному номеру..."
              class="search-input"
          />
          <button v-if="hasActiveFilters" @click="resetFilters" class="btn-icon" title="Сбросить фильтры">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
            </svg>
          </button>
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

    <!-- Equipment Table -->
    <div class="table-container">
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <span>Загрузка данных...</span>
      </div>

      <table v-else class="equipment-table">
        <thead>
        <tr>
          <th @click="sortBy('name')" class="sortable">
            Название
            <svg v-if="sortField === 'name'" class="sort-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline :points="sortDirection === 'asc' ? '6 9 12 15 18 9' : '6 15 12 9 18 15'"/>
            </svg>
          </th>
          <th>Серийный номер</th>
          <th>Категория</th>
          <th @click="sortBy('status')" class="sortable">
            Статус
            <svg v-if="sortField === 'status'" class="sort-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline :points="sortDirection === 'asc' ? '6 9 12 15 18 9' : '6 15 12 9 18 15'"/>
            </svg>
          </th>
          <th @click="sortBy('quantity')" class="sortable">
            Количество
            <svg v-if="sortField === 'quantity'" class="sort-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline :points="sortDirection === 'asc' ? '6 9 12 15 18 9' : '6 15 12 9 18 15'"/>
            </svg>
          </th>
          <th>Местоположение</th>
          <th>Поставщик</th>
          <th @click="sortBy('price')" class="sortable">
            Цена
            <svg v-if="sortField === 'price'" class="sort-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline :points="sortDirection === 'asc' ? '6 9 12 15 18 9' : '6 15 12 9 18 15'"/>
            </svg>
          </th>
          <th @click="sortBy('totalPrice')" class="sortable">
            Общая сумма
            <svg v-if="sortField === 'totalPrice'" class="sort-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline :points="sortDirection === 'asc' ? '6 9 12 15 18 9' : '6 15 12 9 18 15'"/>
            </svg>
          </th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="equipment in paginatedEquipment" 
            :key="equipment.id" 
            class="table-row"
            :class="{ 'selected': selectedEquipment && selectedEquipment.id === equipment.id }"
            @click="selectEquipment(equipment)">
          <td class="equipment-name">
            <div class="name-cell">
              <strong>{{ equipment.name }}</strong>
              <span class="description">{{ equipment.description.substring(0, 50) }}...</span>
            </div>
          </td>
          <td class="serial-number">{{ equipment.serial_number }}</td>
          <td>
            <span class="category-badge">{{ equipment.category?.name || 'Не указано' }}</span>
          </td>
          <td>
              <span :class="['status-badge', `status-${equipment.status}`]">
                {{ getStatusText(equipment.status) }}
              </span>
          </td>
          <td class="quantity">{{ equipment.quantity }}</td>
          <td>{{ equipment.location?.name || 'Не указано' }}</td>
          <td>{{ equipment.supplier?.name || 'Не указано' }}</td>
          <td class="price">{{ formatPrice(equipment.price) }}</td>
          <td class="total-price">{{ formatPrice(equipment.price * equipment.quantity) }}</td>
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

    <!-- Equipment Modal -->
    <div v-if="showModal" class="modal-overlay" @click="closeModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h2>{{ modalMode === 'create' ? 'Добавить оборудование' : modalMode === 'edit' ? 'Редактировать оборудование' : 'Просмотр оборудования' }}</h2>
          <button @click="closeModal" class="modal-close">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <form @submit.prevent="saveEquipment">
            <div class="form-grid">
              <div class="form-group">
                <label>Название *</label>
                <input
                    v-model="currentEquipment.name"
                    :disabled="modalMode === 'view'"
                    required
                    class="form-input"
                />
              </div>

              <div class="form-group">
                <label>Серийный номер *</label>
                <div v-if="modalMode === 'view'" class="form-static-value">
                    {{ currentEquipment.serial_number }}
                </div>
                <div v-else class="serial-number-input">
                    <input
                        v-model="currentEquipment.serial_number"
                        required
                        class="form-input"
                        placeholder="Введите серийный номер"
                    />
                    <button type="button" class="btn btn-secondary" @click="generateSerial">
                        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                            <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
                            <rect x="8" y="2" width="8" height="4" rx="1" ry="1"/>
                            <path d="M12 11v6"/>
                            <path d="M9 14h6"/>
                        </svg>
                        Сгенерировать
                    </button>
                </div>
              </div>

              <div class="form-group">
                <label>Категория *</label>
                <div v-if="modalMode === 'view'" class="form-static-value">
                    {{ currentEquipment.category?.name || '—' }}
                </div>
                <select
                    v-else
                    v-model="currentEquipment.category_id"
                    required
                    class="form-select"
                >
                    <option value="">Выберите категорию</option>
                    <option v-for="category in categories" :key="category.id" :value="category.id">
                        {{ category.name }}
                    </option>
                </select>
              </div>

              <div class="form-group">
                <label>Статус</label>
                <div v-if="modalMode === 'view'" class="form-static-value">
                  <span :class="['status-badge', `status-${currentEquipment.status}`]">
                    {{ getStatusText(currentEquipment.status) }}
                  </span>
                </div>
                <select
                    v-else
                    v-model="currentEquipment.status"
                    class="form-select"
                >
                  <option value="available">Доступно</option>
                  <option value="in_use">Используется</option>
                  <option value="maintenance">На обслуживании</option>
                  <option value="written_off">Списано</option>
                </select>
              </div>

              <div class="form-group">
                <label>Количество *</label>
                <input
                    v-model.number="currentEquipment.quantity"
                    :disabled="modalMode === 'view'"
                    type="number"
                    min="0"
                    required
                    class="form-input"
                />
              </div>

              <div class="form-group">
                <label>Цена</label>
                <input
                    v-model.number="currentEquipment.price"
                    :disabled="modalMode === 'view'"
                    type="number"
                    step="0.01"
                    min="0"
                    class="form-input"
                />
              </div>

              <div class="form-group">
                <label>Местоположение</label>
                <div v-if="modalMode === 'view'" class="form-static-value">
                  {{ currentEquipment.location?.name || 'Не указано' }}
                </div>
                <select
                    v-else
                    v-model="currentEquipment.location_id"
                    class="form-select"
                >
                  <option value="">Выберите местоположение</option>
                  <option v-for="location in locations" :key="location.id" :value="location.id">
                    {{ location.name }}
                  </option>
                </select>
              </div>

              <div class="form-group">
                <label>Поставщик</label>
                <div v-if="modalMode === 'view'" class="form-static-value">
                  {{ currentEquipment.supplier?.name || 'Не указано' }}
                </div>
                <select
                    v-else
                    v-model="currentEquipment.supplier_id"
                    class="form-select"
                >
                  <option value="">Выберите поставщика</option>
                  <option v-for="supplier in suppliers" :key="supplier.id" :value="supplier.id">
                    {{ supplier.name }}
                  </option>
                </select>
              </div>
            </div>

            <div class="form-group full-width">
              <label>Описание</label>
              <textarea
                  v-model="currentEquipment.description"
                  :disabled="modalMode === 'view'"
                  rows="4"
                  class="form-textarea"
              ></textarea>
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

    <!-- Transfer Modal -->
    <div v-if="showTransferModal" class="modal-overlay" @click="closeTransferModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h2>Передача оборудования</h2>
          <button @click="closeTransferModal" class="modal-close">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <form @submit.prevent="transferEquipment">
            <div class="form-group">
              <label>Оборудование</label>
              <div class="form-static-value">{{ currentTransfer.equipment?.name }}</div>
            </div>

            <div class="form-group">
              <label>Текущее местоположение</label>
              <div class="form-static-value">{{ currentTransfer.equipment?.location?.name || 'Не указано' }}</div>
            </div>

            <div class="form-group">
              <label>Новое местоположение *</label>
              <select
                  v-model="currentTransfer.toLocationId"
                  class="form-select"
                  required
              >
                <option value="">Выберите местоположение</option>
                <option v-for="location in locations" 
                        :key="location.id" 
                        :value="location.id"
                        :disabled="location.id === currentTransfer.equipment?.location_id">
                  {{ location.name }}
                </option>
              </select>
            </div>

            <div class="form-group">
              <label>Количество *</label>
              <input
                  v-model.number="currentTransfer.quantity"
                  type="number"
                  min="1"
                  :max="currentTransfer.equipment?.quantity"
                  required
                  class="form-input"
              />
              <small class="form-help">Доступно: {{ currentTransfer.equipment?.quantity }}</small>
            </div>

            <div class="form-group">
              <label>Дата передачи *</label>
              <input
                  v-model="currentTransfer.date"
                  type="date"
                  required
                  class="form-input"
              />
            </div>

            <div class="form-group">
              <label>Причина *</label>
              <select
                  v-model="currentTransfer.reason"
                  class="form-select"
                  required
              >
                <option value="">Выберите причину</option>
                <option value="transfer">Перемещение</option>
                <option value="repair">Ремонт</option>
                <option value="inventory">Инвентаризация</option>
              </select>
            </div>

            <div class="modal-actions">
              <button type="button" @click="closeTransferModal" class="btn btn-secondary">
                Отмена
              </button>
              <button type="submit" class="btn btn-primary">
                Передать
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

<script>
import HeaderComponent from '../components/Header/HeaderComponent.vue'
import {
  CreateEquipment,
  DeleteEquipment,
  GetAllEquipment,
  UpdateEquipment
} from "../../wailsjs/go/service/EquipmentService";
import {
  GetAllLocations,
} from "../../wailsjs/go/service/LocationService";
import {
  GetAllSuppliers,
} from "../../wailsjs/go/service/SupplierService";
import {
  GetAllCategories,
} from "../../wailsjs/go/service/CategoryService";
import {
  CreateMovement,
} from "../../wailsjs/go/service/MovementService";
import { generateSerialNumber } from '../utils/serialNumber'
import { getUser, clearAuth } from '../utils/auth'
import * as XLSX from 'xlsx'
import { utils } from 'xlsx'
import { writeFile } from 'xlsx'

export default {
  name: 'EquipmentView',
  components: {
    HeaderComponent
  },
  data() {
    return {
      equipment: [],
      filteredEquipment: [],
      locations: [],
      suppliers: [],
      categories: [],
      loading: false,
      showModal: false,
      modalMode: 'create', // 'create', 'edit', 'view'
      currentEquipment: this.getEmptyEquipment(),

      // Filters and search
      searchQuery: '',
      categoryFilter: '',
      statusFilter: '',

      // Sorting
      sortField: 'name',
      sortDirection: 'asc',

      // Pagination
      currentPage: 1,
      itemsPerPage: 20,

      // Notification
      notification: {
        show: false,
        message: '',
        type: 'success'
      },

      selectedEquipment: null,
      showTransferModal: false,
      currentTransfer: this.getEmptyTransfer(),
      currentUser: null,
    }
  },

  computed: {
    totalItems() {
      return this.filteredEquipment.length
    },
    totalSum() {
      return this.filteredEquipment.reduce((sum, item) => {
        return sum + (item.price * item.quantity)
      }, 0)
    },

    sortedEquipment() {
      const sorted = [...this.filteredEquipment].sort((a, b) => {
        let aVal = a[this.sortField]
        let bVal = b[this.sortField]

        if (typeof aVal === 'string') {
          aVal = aVal.toLowerCase()
          bVal = bVal.toLowerCase()
        }

        if (this.sortDirection === 'asc') {
          return aVal < bVal ? -1 : aVal > bVal ? 1 : 0
        } else {
          return aVal > bVal ? -1 : aVal < bVal ? 1 : 0
        }
      })

      return sorted
    },

    paginatedEquipment() {
      const start = (this.currentPage - 1) * this.itemsPerPage
      const end = start + this.itemsPerPage
      return this.sortedEquipment.slice(start, end)
    },

    totalPages() {
      return Math.ceil(this.sortedEquipment.length / this.itemsPerPage)
    },

    hasActiveFilters() {
      return this.searchQuery || this.categoryFilter || this.statusFilter
    }
  },

  async mounted() {
    // Проверяем наличие токена
    const token = localStorage.getItem('token')
    if (!token) {
      this.$router.push('/auth')
      return
    }

    await this.loadData()
    // Загружаем текущего пользователя из хранилища
    const user = getUser()
    if (user) {
      this.currentUser = user
    } else {
      clearAuth()
      this.$router.push('/auth')
    }
  },

  methods: {
    async loadData() {
      this.loading = true
      try {
        await Promise.all([
          this.loadEquipment(),
          this.loadLocations(),
          this.loadSuppliers(),
          this.loadCategories()
        ])
      } finally {
        this.loading = false
      }
    },

    async loadEquipment() {
      try {
        const response = await GetAllEquipment()
        if (response.model) {
          this.equipment = response.model
          this.filteredEquipment = [...this.equipment]
        } else {
          this.showNotification(response.msg || 'Ошибка загрузки данных', 'error')
        }
      } catch (error) {
        this.showNotification('Ошибка подключения к серверу', 'error')
      }
    },

    async loadLocations() {
      try {
        const response = await GetAllLocations()
        if (response.model) {
          this.locations = response.model
        }
      } catch (error) {
        console.error('Ошибка загрузки локаций:', error)
      }
    },

    async loadSuppliers() {
      try {
        const response = await GetAllSuppliers()
        if (response.model) {
          this.suppliers = response.model
        }
      } catch (error) {
        console.error('Ошибка загрузки поставщиков:', error)
      }
    },

    async loadCategories() {
      try {
        const response = await GetAllCategories()
        if (response.model) {
          this.categories = response.model
        }
      } catch (error) {
        console.error('Ошибка загрузки категорий:', error)
      }
    },

    applyFilters() {
      let filtered = [...this.equipment]

      // Search filter
      if (this.searchQuery.trim()) {
        const query = this.searchQuery.toLowerCase()
        filtered = filtered.filter(eq =>
            eq.name.toLowerCase().includes(query) ||
            eq.serial_number.toLowerCase().includes(query) ||
            eq.description.toLowerCase().includes(query)
        )
      }

      // Category filter
      if (this.categoryFilter) {
        filtered = filtered.filter(eq => eq.category_id === this.categoryFilter)
      }

      // Status filter
      if (this.statusFilter) {
        filtered = filtered.filter(eq => eq.status === this.statusFilter)
      }

      this.filteredEquipment = filtered
      this.currentPage = 1
    },

    sortBy(field) {
      if (this.sortField === field) {
        this.sortDirection = this.sortDirection === 'asc' ? 'desc' : 'asc'
      } else {
        this.sortField = field
        this.sortDirection = 'asc'
      }
    },

    refreshData() {
      this.loadData()
    },

    // Modal methods
    openCreateModal() {
      this.modalMode = 'create'
      this.currentEquipment = this.getEmptyEquipment()
      this.showModal = true
    },

    viewEquipment(equipment) {
      this.modalMode = 'view'
      this.currentEquipment = { ...equipment }
      this.showModal = true
    },

    editEquipment(equipment) {
      this.modalMode = 'edit'
      this.currentEquipment = { ...equipment }
      this.showModal = true
    },

    closeModal() {
      this.showModal = false
      this.currentEquipment = this.getEmptyEquipment()
    },

    async saveEquipment() {
      try {
        // Проверяем обязательные поля
        if (!this.currentEquipment.name) {
          this.showNotification('Введите название оборудования', 'error')
          return
        }
        if (!this.currentEquipment.serial_number) {
          this.showNotification('Введите серийный номер', 'error')
          return
        }
        if (!this.currentEquipment.category_id) {
          this.showNotification('Выберите категорию', 'error')
          return
        }
        if (this.currentEquipment.quantity <= 0) {
          this.showNotification('Количество должно быть больше нуля', 'error')
          return
        }

        let response
        if (this.modalMode === 'create') {
          response = await CreateEquipment(this.currentEquipment)
        } else {
          response = await UpdateEquipment(this.currentEquipment)
        }

        if (response.model) {
          this.showNotification(response.msg || 'Оборудование успешно сохранено', 'success')
          this.closeModal()
          await this.loadEquipment()
        } else {
          this.showNotification(response.msg || 'Ошибка сохранения', 'error')
        }
      } catch (error) {
        this.showNotification('Ошибка подключения к серверу', 'error')
      }
    },

    async deleteEquipment(equipment) {
      if (!confirm(`Удалить оборудование "${equipment.name}"?`)) {
        return
      }

      try {
        const response = await DeleteEquipment(equipment.id)
        if (response.model) {
          this.showNotification(response.msg, 'success')
          await this.loadEquipment()
        } else {
          this.showNotification(response.msg || 'Ошибка удаления', 'error')
        }
      } catch (error) {
        this.showNotification('Ошибка подключения к серверу', 'error')
      }
    },

    // Utility methods
    async exportToExcel() {
      try {
        const workbook = utils.book_new();
        const worksheet = utils.json_to_sheet(this.filteredEquipment.map(item => ({
          'Название': item.name,
          'Описание': item.description,
          'Серийный номер': item.serial_number,
          'Статус': this.getStatusText(item.status),
          'Количество': item.quantity,
          'Цена': item.price,
          'Общая стоимость': item.price * item.quantity,
          'Категория': item.category?.name || 'Не указано',
          'Местоположение': item.location?.name || 'Не указано',
          'Поставщик': item.supplier?.name || 'Не указано'
        })));

        // Устанавливаем ширину столбцов
        const wscols = [
          { wch: 30 }, // Название
          { wch: 40 }, // Описание
          { wch: 15 }, // Серийный номер
          { wch: 15 }, // Статус
          { wch: 10 }, // Количество
          { wch: 15 }, // Цена
          { wch: 15 }, // Общая стоимость
          { wch: 15 }, // Категория
          { wch: 20 }, // Местоположение
          { wch: 20 }  // Поставщик
        ];
        worksheet['!cols'] = wscols;

        utils.book_append_sheet(workbook, worksheet, 'Оборудование');
        writeFile(workbook, 'Оборудование.xlsx');
        
        this.showNotification('Файл успешно экспортирован', 'success');
      } catch (error) {
        console.error('Ошибка при экспорте в Excel:', error);
        this.showNotification('Ошибка при экспорте в Excel', 'error');
      }
    },

    getEmptyEquipment() {
      return {
        id: 0,
        name: '',
        description: '',
        serial_number: '',
        status: 'available',
        quantity: 1,
        price: 0,
        category_id: '',
        location_id: '',
        supplier_id: ''
      }
    },

    getStatusText(status) {
      const statusMap = {
        'available': 'Доступно',
        'in_use': 'Используется',
        'maintenance': 'На обслуживании',
        'written_off': 'Списано'
      }
      return statusMap[status] || status
    },

    formatPrice(price) {
      if (!price) return '—'
      return new Intl.NumberFormat('ru-RU', {
        style: 'currency',
        currency: 'RUB'
      }).format(price)
    },

    showNotification(message, type = 'success') {
      this.notification = {
        show: true,
        message,
        type
      }

      setTimeout(() => {
        this.notification.show = false
      }, 4000)
    },

    selectEquipment(equipment) {
      this.selectedEquipment = equipment
    },

    getEmptyTransfer() {
      return {
        equipment: null,
        fromLocationId: 0,
        toLocationId: '',
        quantity: 1,
        reason: '',
        createdById: 0,
        date: new Date().toISOString().split('T')[0]
      }
    },

    openTransferModal(equipment) {
      this.currentTransfer = {
        ...this.getEmptyTransfer(),
        equipment: equipment,
        fromLocationId: equipment.location_id,
        quantity: 1
      }
      this.showTransferModal = true
    },

    closeTransferModal() {
      this.showTransferModal = false
      this.currentTransfer = this.getEmptyTransfer()
    },

    async transferEquipment() {
      try {
        if (!this.currentUser) {
          const response = await Login({ username: "admin", password: "admin" })
          if (!response) {
            this.showNotification('Ошибка: пользователь не авторизован', 'error')
            return
          }
          this.currentUser = response
        }

        const movement = {
          equipment_id: this.currentTransfer.equipment.id,
          from_location_id: this.currentTransfer.fromLocationId,
          to_location_id: parseInt(this.currentTransfer.toLocationId),
          quantity: this.currentTransfer.quantity,
          reason: this.currentTransfer.reason,
          created_by_id: this.currentUser.id,
          date: this.currentTransfer.date
        }

        const response = await CreateMovement(movement)
        if (response.model) {
          this.showNotification('Оборудование успешно передано', 'success')
          this.closeTransferModal()
          await this.loadEquipment()
        } else {
          this.showNotification(response.msg || 'Ошибка при передаче оборудования', 'error')
        }
      } catch (error) {
        this.showNotification('Ошибка подключения к серверу', 'error')
      }
    },

    generateSerial() {
      const selectedCategory = this.categories.find(c => c.id === this.currentEquipment.category_id)
      if (selectedCategory) {
        this.currentEquipment.serial_number = generateSerialNumber(selectedCategory.name)
      } else {
        this.showNotification('Сначала выберите категорию', 'error')
      }
    },

    resetFilters() {
      this.searchQuery = ''
      this.categoryFilter = ''
      this.statusFilter = ''
      this.applyFilters()
    }
  }
}
</script>

<style scoped>
.equipment-container {
  padding: 24px;
  background: #f8fafc;
  min-height: 100vh;
  margin-top: 64px;
}

/* Header */
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  position: relative;
  z-index: 1;
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

.btn-primary:hover {
  background: #2563eb;
}

.btn-secondary {
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.btn-secondary:hover {
  background: #e2e8f0;
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

.btn-icon svg {
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
  gap: 8px;
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
  padding: 8px 12px 8px 40px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: white;
  font-size: 14px;
  color: #334155;
  width: 280px;
}

/* Table */
.table-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.equipment-table {
  width: 100%;
  border-collapse: collapse;
}

.equipment-table th {
  background: #f8fafc;
  padding: 16px 12px;
  text-align: left;
  font-weight: 600;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
  font-size: 14px;
}

.equipment-table th.sortable {
  cursor: pointer;
  user-select: none;
  position: relative;
}

.equipment-table th.sortable:hover {
  background: #f1f5f9;
}

.sort-icon {
  width: 14px;
  height: 14px;
  margin-left: 4px;
  vertical-align: middle;
}

.equipment-table td {
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

.equipment-name {
  min-width: 200px;
}

.name-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.description {
  font-size: 12px;
  color: #64748b;
}

.serial-number {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: #475569;
}

.category-badge {
  display: inline-block;
  padding: 4px 8px;
  background: #e0f2fe;
  color: #0369a1;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-available {
  background: #dcfce7;
  color: #166534;
}

.status-in_use {
  background: #fef3c7;
  color: #92400e;
}

.status-maintenance {
  background: #fed7d7;
  color: #c53030;
}

.status-written_off {
  background: #f3f4f6;
  color: #4b5563;
}

.quantity {
  font-weight: 600;
  color: #3b82f6;
}

.price {
  font-weight: 600;
  color: #059669;
}

.total-price {
  font-weight: 600;
  color: #059669;
}

.actions-col {
  width: 120px;
}

.action-buttons {
  display: flex;
  gap: 4px;
}

/* Loading */
.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #64748b;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f1f5f9;
  border-top: 3px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Pagination */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid #f1f5f9;
}

.pagination-btn {
  padding: 8px 12px;
  border: 1px solid #e2e8f0;
  background: white;
  color: #475569;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.pagination-btn:hover:not(:disabled) {
  background: #f8fafc;
  border-color: #cbd5e1;
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-info {
  font-size: 14px;
  color: #64748b;
  margin: 0 8px;
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
  max-width: 800px;
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

.modal-close svg {
  width: 16px;
  height: 16px;
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

.form-group.full-width {
  grid-column: 1 / -1;
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
  .equipment-container {
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

  .equipment-table {
    font-size: 12px;
  }

  .equipment-table th,
  .equipment-table td {
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

.btn-info {
  background: #0ea5e9;
  color: white;
}

.btn-info:hover:not(:disabled) {
  background: #0284c7;
}

.form-help {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
}

.serial-number-input {
  display: flex;
  gap: 8px;
}

.serial-number-input .form-input {
  flex: 1;
}
</style>
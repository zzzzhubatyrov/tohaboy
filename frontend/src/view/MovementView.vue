<template>
  <HeaderComponent />
  <div class="movement-container">
    <!-- Header -->
    <div class="header">
      <h1 class="title">Перемещение оборудования</h1>
      <div class="header-stats">
        <div class="stat-item">
          <span class="stat-number">{{ totalMovements }}</span>
          <span class="stat-label">Всего перемещений</span>
        </div>
        <div class="stat-item">
          <span class="stat-number">{{ todayMovements }}</span>
          <span class="stat-label">За сегодня</span>
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
            Создать перемещение
          </button>
        </div>

        <div class="filter-group">
          <select v-model="reasonFilter" @change="applyFilters" class="select">
            <option value="">Все причины</option>
            <option value="transfer">Перемещение</option>
            <option value="inventory">Инвентаризация</option>
            <option value="repair">Ремонт</option>
          </select>

          <select v-model="locationFilter" @change="applyFilters" class="select">
            <option value="">Все местоположения</option>
            <option v-for="location in locations" :key="location.id" :value="location.id">
              {{ location.name }}
            </option>
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
              placeholder="Поиск по оборудованию..."
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

    <!-- Movements Table -->
    <div class="table-container">
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <span>Загрузка данных...</span>
      </div>

      <table v-else class="movement-table">
        <thead>
        <tr>
          <th @click="sortBy('date')" class="sortable">
            Дата
            <svg v-if="sortField === 'date'" class="sort-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline :points="sortDirection === 'asc' ? '6 9 12 15 18 9' : '6 15 12 9 18 15'"/>
            </svg>
          </th>
          <th>Оборудование</th>
          <th>Откуда</th>
          <th>Куда</th>
          <th>Количество</th>
          <th>Причина</th>
          <th>Документ</th>
          <th>Пользователь</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="movement in paginatedMovements" 
            :key="movement.id" 
            class="table-row"
            @click="viewMovement(movement)">
          <td>{{ formatDate(movement.date) }}</td>
          <td>
            <div class="equipment-cell">
              <strong>{{ movement.equipment?.name }}</strong>
              <span class="serial-number">{{ movement.equipment?.serial_number }}</span>
            </div>
          </td>
          <td>{{ getLocationName(movement.from_location_id) || '—' }}</td>
          <td>{{ getLocationName(movement.to_location_id) }}</td>
          <td class="quantity">{{ movement.quantity }}</td>
          <td>
            <span :class="['reason-badge', `reason-${movement.reason}`]">
              {{ getReasonText(movement.reason) }}
            </span>
          </td>
          <td>{{ movement.document?.number || '—' }}</td>
          <td>{{ movement.user?.username }}</td>
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

    <!-- Movement Modal -->
    <div v-if="showModal" class="modal-overlay" @click="closeModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h2>{{ modalMode === 'create' ? 'Создать перемещение' : 'Просмотр перемещения' }}</h2>
          <button @click="closeModal" class="modal-close">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <form @submit.prevent="saveMovement">
            <div class="form-grid">
              <div class="form-group">
                <label>Оборудование *</label>
                <div v-if="modalMode === 'view'" class="form-static-value">
                  {{ currentMovement.equipment?.name }}
                </div>
                <select
                    v-else
                    v-model="currentMovement.equipment_id"
                    class="form-select"
                    required
                >
                  <option value="">Выберите оборудование</option>
                  <option v-for="equipment in availableEquipment" :key="equipment.id" :value="equipment.id">
                    {{ equipment.name }} ({{ equipment.serial_number }})
                  </option>
                </select>
              </div>

              <div class="form-group">
                <label>Откуда</label>
                <div v-if="modalMode === 'view'" class="form-static-value">
                  {{ getLocationName(currentMovement.from_location_id) || '—' }}
                </div>
                <select
                    v-else
                    v-model="currentMovement.from_location_id"
                    class="form-select"
                >
                  <option value="">Выберите местоположение</option>
                  <option v-for="location in locations" :key="location.id" :value="location.id">
                    {{ location.name }}
                  </option>
                </select>
              </div>

              <div class="form-group">
                <label>Куда *</label>
                <div v-if="modalMode === 'view'" class="form-static-value">
                  {{ getLocationName(currentMovement.to_location_id) }}
                </div>
                <select
                    v-else
                    v-model="currentMovement.to_location_id"
                    class="form-select"
                    required
                >
                  <option value="">Выберите местоположение</option>
                  <option v-for="location in locations" :key="location.id" :value="location.id">
                    {{ location.name }}
                  </option>
                </select>
              </div>

              <div class="form-group">
                <label>Количество *</label>
                <div v-if="modalMode === 'view'" class="form-static-value">
                  {{ currentMovement.quantity }}
                </div>
                <input
                    v-else
                    v-model.number="currentMovement.quantity"
                    type="number"
                    min="1"
                    :max="maxQuantity"
                    required
                    class="form-input"
                />
              </div>

              <div class="form-group">
                <label>Причина *</label>
                <div v-if="modalMode === 'view'" class="form-static-value">
                  <span :class="['reason-badge', `reason-${currentMovement.reason}`]">
                    {{ getReasonText(currentMovement.reason) }}
                  </span>
                </div>
                <select
                    v-else
                    v-model="currentMovement.reason"
                    class="form-select"
                    required
                >
                  <option value="">Выберите причину</option>
                  <option value="transfer">Перемещение</option>
                  <option value="inventory">Инвентаризация</option>
                  <option value="repair">Ремонт</option>
                </select>
              </div>

              <div class="form-group">
                <label>Документ</label>
                <div v-if="modalMode === 'view'" class="form-static-value">
                  {{ currentMovement.document?.number || '—' }}
                </div>
                <select
                    v-else
                    v-model="currentMovement.document_id"
                    class="form-select"
                >
                  <option value="">Выберите документ</option>
                  <option v-for="doc in availableDocuments" :key="doc.id" :value="doc.id">
                    {{ doc.number }}
                  </option>
                </select>
              </div>
            </div>

            <div v-if="modalMode !== 'view'" class="modal-actions">
              <button type="button" @click="closeModal" class="btn btn-secondary">
                Отмена
              </button>
              <button type="submit" class="btn btn-primary">
                Создать
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
  CreateMovement,
  GetAllMovements,
  GetAllLocations,
  GetAllEquipment,
  GetAllDocuments
} from "../../wailsjs/go/service/MovementService.js";

export default {
  name: 'MovementView',
  components: {
    HeaderComponent
  },
  data() {
    return {
      movements: [],
      filteredMovements: [],
      locations: [],
      equipment: [],
      documents: [],
      loading: false,
      showModal: false,
      modalMode: 'create', // 'create' или 'view'
      currentMovement: this.getEmptyMovement(),

      // Фильтры и поиск
      searchQuery: '',
      reasonFilter: '',
      locationFilter: '',

      // Сортировка
      sortField: 'date',
      sortDirection: 'desc',

      // Пагинация
      currentPage: 1,
      itemsPerPage: 20,

      // Уведомления
      notification: {
        show: false,
        message: '',
        type: 'success'
      }
    }
  },

  computed: {
    totalMovements() {
      return this.movements.length
    },

    todayMovements() {
      const today = new Date()
      today.setHours(0, 0, 0, 0)
      return this.movements.filter(m => new Date(m.date) >= today).length
    },

    availableEquipment() {
      return this.equipment.filter(e => e.status !== 'written_off' && e.quantity > 0)
    },

    availableDocuments() {
      return this.documents.filter(d => d.status === 'draft')
    },

    maxQuantity() {
      if (!this.currentMovement.equipment_id) return 1
      const equipment = this.equipment.find(e => e.id === this.currentMovement.equipment_id)
      return equipment ? equipment.quantity : 1
    },

    sortedMovements() {
      return [...this.filteredMovements].sort((a, b) => {
        let aVal = a[this.sortField]
        let bVal = b[this.sortField]

        if (this.sortField === 'date') {
          aVal = new Date(aVal)
          bVal = new Date(bVal)
        }

        if (this.sortDirection === 'asc') {
          return aVal < bVal ? -1 : aVal > bVal ? 1 : 0
        } else {
          return aVal > bVal ? -1 : aVal < bVal ? 1 : 0
        }
      })
    },

    paginatedMovements() {
      const start = (this.currentPage - 1) * this.itemsPerPage
      const end = start + this.itemsPerPage
      return this.sortedMovements.slice(start, end)
    },

    totalPages() {
      return Math.ceil(this.sortedMovements.length / this.itemsPerPage)
    }
  },

  async mounted() {
    await this.loadData()
  },

  methods: {
    async loadData() {
      this.loading = true
      try {
        await Promise.all([
          this.loadMovements(),
          this.loadLocations(),
          this.loadEquipment(),
          this.loadDocuments()
        ])
      } finally {
        this.loading = false
      }
    },

    async loadMovements() {
      try {
        const response = await GetAllMovements()
        if (response.model) {
          this.movements = response.model
          this.filteredMovements = [...this.movements]
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
        console.error('Ошибка загрузки местоположений:', error)
      }
    },

    async loadEquipment() {
      try {
        const response = await GetAllEquipment()
        if (response.model) {
          this.equipment = response.model
        }
      } catch (error) {
        console.error('Ошибка загрузки оборудования:', error)
      }
    },

    async loadDocuments() {
      try {
        const response = await GetAllDocuments()
        if (response.model) {
          this.documents = response.model
        }
      } catch (error) {
        console.error('Ошибка загрузки документов:', error)
      }
    },

    applyFilters() {
      let filtered = [...this.movements]

      // Поиск
      if (this.searchQuery.trim()) {
        const query = this.searchQuery.toLowerCase()
        filtered = filtered.filter(m =>
            m.equipment?.name.toLowerCase().includes(query) ||
            m.equipment?.serial_number.toLowerCase().includes(query)
        )
      }

      // Фильтр по причине
      if (this.reasonFilter) {
        filtered = filtered.filter(m => m.reason === this.reasonFilter)
      }

      // Фильтр по местоположению
      if (this.locationFilter) {
        const locationId = parseInt(this.locationFilter)
        filtered = filtered.filter(m =>
            m.from_location_id === locationId || m.to_location_id === locationId
        )
      }

      this.filteredMovements = filtered
      this.currentPage = 1
    },

    sortBy(field) {
      if (this.sortField === field) {
        this.sortDirection = this.sortDirection === 'asc' ? 'desc' : 'asc'
      } else {
        this.sortField = field
        this.sortDirection = 'desc'
      }
    },

    refreshData() {
      this.loadData()
    },

    openCreateModal() {
      this.modalMode = 'create'
      this.currentMovement = this.getEmptyMovement()
      this.showModal = true
    },

    viewMovement(movement) {
      this.modalMode = 'view'
      this.currentMovement = { ...movement }
      this.showModal = true
    },

    closeModal() {
      this.showModal = false
      this.currentMovement = this.getEmptyMovement()
    },

    async saveMovement() {
      try {
        const response = await CreateMovement(this.currentMovement)
        if (response.model) {
          this.showNotification('Перемещение создано успешно', 'success')
          this.closeModal()
          await this.loadData()
        } else {
          this.showNotification(response.msg || 'Ошибка создания перемещения', 'error')
        }
      } catch (error) {
        this.showNotification('Ошибка подключения к серверу', 'error')
      }
    },

    getEmptyMovement() {
      return {
        id: 0,
        equipment_id: '',
        from_location_id: '',
        to_location_id: '',
        quantity: 1,
        reason: '',
        document_id: '',
        date: new Date().toISOString()
      }
    },

    getLocationName(id) {
      if (!id) return null
      const location = this.locations.find(l => l.id === id)
      return location ? location.name : null
    },

    getReasonText(reason) {
      const reasonMap = {
        'transfer': 'Перемещение',
        'inventory': 'Инвентаризация',
        'repair': 'Ремонт'
      }
      return reasonMap[reason] || reason
    },

    formatDate(date) {
      return new Date(date).toLocaleString('ru-RU', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
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
    }
  }
}
</script>

<style scoped>
.movement-container {
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
}

.filter-group {
  display: flex;
  gap: 12px;
}

/* Table */
.table-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.movement-table {
  width: 100%;
  border-collapse: collapse;
}

.movement-table th {
  background: #f8fafc;
  padding: 16px 12px;
  text-align: left;
  font-weight: 600;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
  font-size: 14px;
}

.movement-table td {
  padding: 16px 12px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  font-size: 14px;
}

.equipment-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.serial-number {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #64748b;
}

.quantity {
  font-weight: 600;
  color: #3b82f6;
}

.reason-badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.reason-transfer {
  background: #e0f2fe;
  color: #0369a1;
}

.reason-inventory {
  background: #fef3c7;
  color: #92400e;
}

.reason-repair {
  background: #fed7d7;
  color: #c53030;
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

.modal-body {
  padding: 24px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
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

/* Responsive */
@media (max-width: 1200px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .movement-container {
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
  }

  .filter-group {
    flex-direction: column;
  }

  .modal {
    width: 95%;
    margin: 10px;
  }
}
</style> 
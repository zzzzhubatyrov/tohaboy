<template>
  <HeaderComponent />
  <div class="movements-container">
    <div class="header">
      <h1 class="title">История перемещений</h1>
    </div>

    <div class="filters-panel">
      <div class="filters-left">
        <select v-model="reasonFilter" @change="applyFilters" class="select">
          <option value="">Все причины</option>
          <option value="transfer">Перемещение</option>
          <option value="repair">Ремонт</option>
          <option value="inventory">Инвентаризация</option>
        </select>
      </div>

      <div class="filters-right">
        <div class="search-box">
          <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <circle cx="11" cy="11" r="8"/>
            <path d="m21 21-4.35-4.35"/>
          </svg>
          <input
              v-model="searchQuery"
              @input="applyFilters"
              placeholder="Поиск по названию оборудования..."
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

    <div class="table-container">
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <span>Загрузка данных...</span>
      </div>

      <table v-else class="movements-table">
        <thead>
        <tr>
          <th @click="sortBy('created_at')" class="sortable">
            Дата
            <svg v-if="sortField === 'created_at'" class="sort-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline :points="sortDirection === 'asc' ? '6 9 12 15 18 9' : '6 15 12 9 18 15'"/>
            </svg>
          </th>
          <th>Оборудование</th>
          <th>Откуда</th>
          <th>Куда</th>
          <th @click="sortBy('quantity')" class="sortable">
            Количество
            <svg v-if="sortField === 'quantity'" class="sort-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline :points="sortDirection === 'asc' ? '6 9 12 15 18 9' : '6 15 12 9 18 15'"/>
            </svg>
          </th>
          <th>Причина</th>
          <th>Кем создано</th>
          <th>Документ</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="movement in paginatedMovements" :key="movement.id" class="table-row">
          <td>{{ formatDate(movement.created_at) }}</td>
          <td class="equipment-name">
            <div class="name-cell">
              <strong>{{ movement.equipment?.name }}</strong>
              <span class="serial-number">{{ movement.equipment?.serial_number }}</span>
            </div>
          </td>
          <td>{{ movement.from_location?.name || 'Не указано' }}</td>
          <td>{{ movement.to_location?.name || 'Не указано' }}</td>
          <td class="quantity">{{ movement.quantity }}</td>
          <td>
            <span :class="['reason-badge', `reason-${movement.reason}`]">
              {{ getReasonText(movement.reason) }}
            </span>
          </td>
          <td>{{ movement.created_by?.username || 'Система' }}</td>
          <td>
            <button v-if="movement.document_id" 
                    @click="viewDocument(movement.document_id)"
                    class="btn btn-link">
              Просмотр
            </button>
            <span v-else>—</span>
          </td>
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

    <!-- Notification -->
    <div v-if="notification.show" :class="['notification', `notification-${notification.type}`]">
      {{ notification.message }}
    </div>
  </div>
</template>

<script>
import HeaderComponent from '../components/Header/HeaderComponent.vue'
import {
  GetAllMovements,
} from "../../wailsjs/go/service/MovementService";
import {
  GetAllEquipment,
} from "../../wailsjs/go/service/EquipmentService";
import {
  GetAllLocations,
} from "../../wailsjs/go/service/LocationService";
import {
  GetDocument,
} from "../../wailsjs/go/service/DocumentService";

export default {
  name: 'MovementsView',
  components: {
    HeaderComponent
  },
  data() {
    return {
      movements: [],
      filteredMovements: [],
      equipmentList: [],
      locations: [],
      loading: false,

      // Filters
      searchQuery: '',
      reasonFilter: '',

      // Sorting
      sortField: 'created_at',
      sortDirection: 'desc',

      // Pagination
      currentPage: 1,
      itemsPerPage: 20,

      // Notification
      notification: {
        show: false,
        message: '',
        type: 'success'
      },
    }
  },

  computed: {
    sortedMovements() {
      const sorted = [...this.filteredMovements].sort((a, b) => {
        let aVal = a[this.sortField]
        let bVal = b[this.sortField]

        if (this.sortField === 'created_at') {
          aVal = new Date(aVal)
          bVal = new Date(bVal)
        }

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
          this.loadEquipment(),
          this.loadLocations()
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

    async loadEquipment() {
      try {
        const response = await GetAllEquipment()
        if (response.model) {
          this.equipmentList = response.model
        }
      } catch (error) {
        console.error('Ошибка загрузки оборудования:', error)
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

    applyFilters() {
      let filtered = [...this.movements]

      // Search filter
      if (this.searchQuery.trim()) {
        const query = this.searchQuery.toLowerCase()
        filtered = filtered.filter(movement =>
            movement.equipment?.name.toLowerCase().includes(query) ||
            movement.equipment?.serial_number.toLowerCase().includes(query)
        )
      }

      // Reason filter
      if (this.reasonFilter) {
        filtered = filtered.filter(movement => movement.reason === this.reasonFilter)
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

    getLocationName(locationId) {
      const location = this.locations.find(l => l.id === locationId)
      return location ? location.name : '—'
    },

    getReasonText(reason) {
      const reasonMap = {
        'transfer': 'Перемещение',
        'repair': 'Ремонт',
        'inventory': 'Инвентаризация'
      }
      return reasonMap[reason] || reason
    },

    formatDate(date) {
      if (!date) return '—'
      return new Date(date).toLocaleString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    },

    async viewDocument(documentId) {
      try {
        const response = await GetDocument(documentId)
        if (response.model) {
          // Здесь можно добавить логику просмотра документа
          // Например, открыть модальное окно или перейти на страницу документа
          this.$router.push(`/documents?id=${documentId}`)
        } else {
          this.showNotification('Ошибка загрузки документа', 'error')
        }
      } catch (error) {
        this.showNotification('Ошибка подключения к серверу', 'error')
      }
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
  }
}
</script>

<style scoped>
.movements-container {
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

.filters-panel {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.filters-left {
  display: flex;
  gap: 12px;
}

.filters-right {
  display: flex;
  gap: 12px;
  align-items: center;
}

.select {
  padding: 8px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: white;
  font-size: 14px;
  color: #334155;
  min-width: 180px;
}

.select:focus {
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
  padding: 8px 12px 8px 40px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  background: white;
  font-size: 14px;
  color: #334155;
  width: 280px;
}

.search-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.table-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.movements-table {
  width: 100%;
  border-collapse: collapse;
}

.movements-table th {
  background: #f8fafc;
  padding: 16px 12px;
  text-align: left;
  font-weight: 600;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
  font-size: 14px;
}

.movements-table th.sortable {
  cursor: pointer;
  user-select: none;
  position: relative;
}

.movements-table th.sortable:hover {
  background: #f1f5f9;
}

.sort-icon {
  width: 14px;
  height: 14px;
  margin-left: 4px;
  vertical-align: middle;
}

.movements-table td {
  padding: 16px 12px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  font-size: 14px;
}

.table-row:hover {
  background: #f8fafc;
}

.equipment-name {
  min-width: 200px;
}

.name-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.serial-number {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: #475569;
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

.reason-repair {
  background: #fef3c7;
  color: #92400e;
}

.reason-inventory {
  background: #dcfce7;
  color: #166534;
}

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

.btn-secondary {
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.btn-secondary:hover {
  background: #e2e8f0;
}

.btn-link {
  background: none;
  color: #3b82f6;
  padding: 4px 8px;
}

.btn-link:hover {
  background: #f1f5f9;
}

.icon {
  width: 16px;
  height: 16px;
  stroke-width: 2;
}

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

@media (max-width: 1200px) {
  .filters-panel {
    flex-direction: column;
    gap: 16px;
  }

  .filters-left, .filters-right {
    width: 100%;
  }

  .search-input {
    width: 100%;
  }
}

@media (max-width: 768px) {
  .movements-container {
    padding: 16px;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .movements-table {
    font-size: 12px;
  }

  .movements-table th,
  .movements-table td {
    padding: 8px 6px;
  }
}
</style> 
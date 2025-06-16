<template>
  <HeaderComponent />
  <div class="suppliers-container">
    <!-- Header -->
    <div class="header">
      <h1 class="title">Управление поставщиками</h1>
      <div class="header-stats">
        <div class="stat-item">
          <span class="stat-number">{{ suppliers.length }}</span>
          <span class="stat-label">Всего поставщиков</span>
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
            Добавить поставщика
          </button>

          <button class="btn btn-secondary" @click="viewSupplier(selectedSupplier)" :disabled="!selectedSupplier">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
              <circle cx="12" cy="12" r="3"/>
            </svg>
            Просмотр
          </button>

          <button class="btn btn-secondary" @click="editSupplier(selectedSupplier)" :disabled="!selectedSupplier">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="m18.5 2.5 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
            Редактировать
          </button>

          <button class="btn btn-danger" @click="deleteSupplier(selectedSupplier)" :disabled="!selectedSupplier">
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline points="3,6 5,6 21,6"/>
              <path d="m19,6v14a2,2 0 0,1 -2,2H7a2,2 0 0,1 -2,-2V6m3,0V4a2,2 0 0,1 2,-2h4a2,2 0 0,1 2,2v2"/>
            </svg>
            Удалить
          </button>
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
              placeholder="Поиск по всем полям..."
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

    <!-- Suppliers Table -->
    <div class="table-container">
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <span>Загрузка данных...</span>
      </div>

      <table v-else class="supplier-table">
        <thead>
        <tr>
          <th @click="sortBy('name')" class="sortable">
            Название
            <svg v-if="sortField === 'name'" class="sort-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <polyline :points="sortDirection === 'asc' ? '6 9 12 15 18 9' : '6 15 12 9 18 15'"/>
            </svg>
          </th>
          <th>Описание</th>
          <th>Контактная информация</th>
          <th>Количество оборудования</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="supplier in paginatedSuppliers"
            :key="supplier.id"
            class="table-row"
            :class="{ 'selected': selectedSupplier && selectedSupplier.id === supplier.id }"
            @click="selectSupplier(supplier)">
          <td class="supplier-name">{{ supplier.name }}</td>
          <td class="supplier-description">{{ supplier.description }}</td>
          <td class="contact-info">
            {{ supplier.phone }}<br>
            {{ supplier.address }}
          </td>
          <td class="equipment-count">{{ supplier.equipment?.length || 0 }}</td>
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

    <!-- Supplier Modal -->
    <div v-if="showModal" class="modal-overlay" @click="closeModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h2>{{ modalMode === 'create' ? 'Добавить поставщика' : modalMode === 'edit' ? 'Редактировать поставщика' : 'Просмотр поставщика' }}</h2>
          <button @click="closeModal" class="modal-close">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <form @submit.prevent="saveSupplier">
            <div class="form-group">
              <label>Название *</label>
              <input
                  v-model="currentSupplier.name"
                  :disabled="modalMode === 'view'"
                  required
                  class="form-input"
              />
            </div>

            <div class="form-group">
              <label>Описание</label>
              <textarea
                  v-model="currentSupplier.description"
                  :disabled="modalMode === 'view'"
                  rows="3"
                  class="form-textarea"
                  placeholder="Описание компании..."
              ></textarea>
            </div>

            <div class="form-group">
              <label>Телефон</label>
              <input
                  v-model="currentSupplier.phone"
                  :disabled="modalMode === 'view'"
                  class="form-input"
                  placeholder="+7 (XXX) XXX-XX-XX"
              />
            </div>

            <div class="form-group">
              <label>Адрес</label>
              <textarea
                  v-model="currentSupplier.address"
                  :disabled="modalMode === 'view'"
                  rows="2"
                  class="form-textarea"
                  placeholder="Физический адрес компании..."
              ></textarea>
            </div>

            <div v-if="modalMode === 'view' && currentSupplier.equipment?.length > 0" class="form-group">
              <label>Поставляемое оборудование</label>
              <div class="equipment-list">
                <div v-for="item in currentSupplier.equipment" :key="item.id" class="equipment-item">
                  <span class="equipment-name">{{ item.name }}</span>
                  <span class="equipment-serial">{{ item.serial_number }}</span>
                </div>
              </div>
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

<script>
import {
  CreateSupplier,
  DeleteSupplier,
  GetAllSuppliers,
  UpdateSupplier
} from "../../wailsjs/go/service/SupplierService.js";
import HeaderComponent from "../components/Header/HeaderComponent.vue";

export default {
  name: 'SupplierView',
  components: {
    HeaderComponent
  },
  data() {
    return {
      suppliers: [],
      filteredSuppliers: [],
      loading: false,
      showModal: false,
      modalMode: 'create', // 'create', 'edit', 'view'
      currentSupplier: this.getEmptySupplier(),
      selectedSupplier: null,

      // Search
      searchQuery: '',

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
      }
    }
  },

  computed: {
    sortedSuppliers() {
      return [...this.filteredSuppliers].sort((a, b) => {
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
    },

    paginatedSuppliers() {
      const start = (this.currentPage - 1) * this.itemsPerPage
      const end = start + this.itemsPerPage
      return this.sortedSuppliers.slice(start, end)
    },

    totalPages() {
      return Math.ceil(this.sortedSuppliers.length / this.itemsPerPage)
    }
  },

  async mounted() {
    await this.loadData()
  },

  methods: {
    async loadData() {
      this.loading = true
      try {
        const response = await GetAllSuppliers()
        if (response.model) {
          this.suppliers = response.model
          this.filteredSuppliers = [...this.suppliers]
        } else {
          this.showNotification(response.msg || 'Ошибка загрузки данных', 'error')
        }
      } catch (error) {
        this.showNotification('Ошибка подключения к серверу', 'error')
      } finally {
        this.loading = false
      }
    },

    applyFilters() {
      if (this.searchQuery.trim()) {
        const query = this.searchQuery.toLowerCase()
        this.filteredSuppliers = this.suppliers.filter(supplier => {
          return (
            supplier.name?.toLowerCase().includes(query) ||
            supplier.description?.toLowerCase().includes(query) ||
            supplier.phone?.toLowerCase().includes(query) ||
            supplier.address?.toLowerCase().includes(query) ||
            supplier.equipment?.some(eq => 
              eq.name?.toLowerCase().includes(query) ||
              eq.serial_number?.toLowerCase().includes(query) ||
              eq.description?.toLowerCase().includes(query)
            )
          )
        })
      } else {
        this.filteredSuppliers = [...this.suppliers]
      }
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

    selectSupplier(supplier) {
      this.selectedSupplier = supplier
    },

    // Modal methods
    openCreateModal() {
      this.modalMode = 'create'
      this.currentSupplier = this.getEmptySupplier()
      this.showModal = true
    },

    viewSupplier(supplier) {
      this.modalMode = 'view'
      this.currentSupplier = { ...supplier }
      this.showModal = true
    },

    editSupplier(supplier) {
      this.modalMode = 'edit'
      this.currentSupplier = { ...supplier }
      this.showModal = true
    },

    closeModal() {
      this.showModal = false
      this.currentSupplier = this.getEmptySupplier()
    },

    async saveSupplier() {
      try {
        let response

        if (this.modalMode === 'create') {
          response = await CreateSupplier(this.currentSupplier)
        } else {
          response = await UpdateSupplier(this.currentSupplier)
        }

        if (response.model) {
          this.showNotification(response.msg, 'success')
          this.closeModal()
          await this.loadData()
        } else {
          this.showNotification(response.msg || 'Ошибка сохранения', 'error')
        }
      } catch (error) {
        this.showNotification('Ошибка подключения к серверу', 'error')
      }
    },

    async deleteSupplier(supplier) {
      if (!confirm(`Удалить поставщика "${supplier.name}"?`)) {
        return
      }

      try {
        const response = await DeleteSupplier(supplier.id)
        if (response.model) {
          this.showNotification(response.msg, 'success')
          await this.loadData()
        } else {
          this.showNotification(response.msg || 'Ошибка удаления', 'error')
        }
      } catch (error) {
        this.showNotification('Ошибка подключения к серверу', 'error')
      }
    },

    getEmptySupplier() {
      return {
        id: 0,
        name: '',
        description: '',
        address: '',
        phone: '',
        equipment: []
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
    }
  }
}
</script>

<style scoped>
.suppliers-container {
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

.action-right {
  display: flex;
  align-items: center;
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

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

.btn-danger {
  background: #ef4444;
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background: #dc2626;
}

.icon {
  width: 16px;
  height: 16px;
  stroke-width: 2;
}

/* Search */
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

/* Table */
.table-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.supplier-table {
  width: 100%;
  border-collapse: collapse;
}

.supplier-table th {
  background: #f8fafc;
  padding: 16px 12px;
  text-align: left;
  font-weight: 600;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
  font-size: 14px;
}

.supplier-table th.sortable {
  cursor: pointer;
  user-select: none;
}

.supplier-table th.sortable:hover {
  background: #f1f5f9;
}

.sort-icon {
  width: 14px;
  height: 14px;
  margin-left: 4px;
  vertical-align: middle;
}

.supplier-table td {
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

.supplier-name {
  font-weight: 600;
}

.supplier-description {
  color: #64748b;
  white-space: pre-line;
}

.contact-info {
  color: #64748b;
  white-space: pre-line;
}

.equipment-count {
  font-weight: 600;
  color: #3b82f6;
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
  max-width: 600px;
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

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
  font-size: 14px;
}

.form-input, .form-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  color: #374151;
  background: white;
  transition: all 0.2s;
}

.form-input:focus, .form-textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-input:disabled, .form-textarea:disabled {
  background: #f9fafb;
  color: #6b7280;
  cursor: not-allowed;
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.equipment-list {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 12px;
  max-height: 200px;
  overflow-y: auto;
}

.equipment-item {
  display: flex;
  justify-content: space-between;
  padding: 8px;
  border-bottom: 1px solid #e2e8f0;
}

.equipment-item:last-child {
  border-bottom: none;
}

.equipment-name {
  font-weight: 500;
  color: #334155;
}

.equipment-serial {
  color: #64748b;
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #f1f5f9;
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
@media (max-width: 768px) {
  .suppliers-container {
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

  .action-right {
    width: 100%;
  }

  .search-input {
    width: 100%;
  }

  .action-buttons-group {
    width: 100%;
  }

  .btn {
    flex: 1;
  }

  .modal {
    width: 95%;
    margin: 10px;
  }
}
</style> 
<template>
  <HeaderComponent />
  <div class="employees-container">
    <div class="header">
      <h1 class="title">Сотрудники</h1>
    </div>

    <div class="action-panel">
      <div class="action-left">
        <button class="btn btn-primary" @click="openCreateModal">
          <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <path d="M12 5v14M5 12h14"/>
          </svg>
          Добавить сотрудника
        </button>
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
              placeholder="Поиск по имени или должности..."
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

      <table v-else class="employees-table">
        <thead>
          <tr>
            <th>ФИО</th>
            <th>Должность</th>
            <th>Отдел</th>
            <th>Контакты</th>
            <th>Оборудование</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="employee in employees" :key="employee.id" class="table-row">
            <td>{{ employee.name }}</td>
            <td>{{ employee.position }}</td>
            <td>{{ employee.department }}</td>
            <td>{{ employee.contact }}</td>
            <td>
              <button class="btn-link" @click="viewEquipment(employee)">
                Просмотреть ({{ employee.equipment?.length || 0 }})
              </button>
            </td>
            <td>
              <div class="actions">
                <button @click="editEmployee(employee)" class="btn-icon" title="Редактировать">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                  </svg>
                </button>
                <button @click="deleteEmployee(employee.id)" class="btn-icon" title="Удалить">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                    <path d="M3 6h18"/>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Employee Modal -->
    <div v-if="showModal" class="modal-overlay" @click="closeModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h2>{{ modalMode === 'create' ? 'Добавить сотрудника' : 'Редактировать сотрудника' }}</h2>
          <button @click="closeModal" class="modal-close">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <form @submit.prevent="saveEmployee">
            <div class="form-group">
              <label>ФИО *</label>
              <input
                  v-model="currentEmployee.name"
                  required
                  class="form-input"
                  placeholder="Введите ФИО сотрудника"
              />
            </div>

            <div class="form-group">
              <label>Должность *</label>
              <input
                  v-model="currentEmployee.position"
                  required
                  class="form-input"
                  placeholder="Введите должность"
              />
            </div>

            <div class="form-group">
              <label>Отдел *</label>
              <input
                  v-model="currentEmployee.department"
                  required
                  class="form-input"
                  placeholder="Введите отдел"
              />
            </div>

            <div class="form-group">
              <label>Контакты</label>
              <input
                  v-model="currentEmployee.contact"
                  class="form-input"
                  placeholder="Введите контактные данные"
              />
            </div>

            <div class="modal-actions">
              <button type="button" @click="closeModal" class="btn btn-secondary">
                Отмена
              </button>
              <button type="submit" class="btn btn-primary">
                {{ modalMode === 'create' ? 'Добавить' : 'Сохранить' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Equipment Modal -->
    <div v-if="showEquipmentModal" class="modal-overlay" @click="closeEquipmentModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h2>Оборудование сотрудника: {{ selectedEmployee?.name }}</h2>
          <button @click="closeEquipmentModal" class="modal-close">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <table class="equipment-table">
            <thead>
              <tr>
                <th>Название</th>
                <th>Серийный номер</th>
                <th>Категория</th>
                <th>Количество</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in selectedEmployee?.equipment" :key="item.id">
                <td>{{ item.name }}</td>
                <td>{{ item.serial_number }}</td>
                <td>{{ item.category?.name }}</td>
                <td>{{ item.quantity }}</td>
              </tr>
            </tbody>
          </table>
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
  CreateEmployee,
  GetAllEmployees,
  UpdateEmployee,
  DeleteEmployee,
} from "../../wailsjs/go/service/EmployeeService"

export default {
  name: 'EmployeesView',
  components: {
    HeaderComponent
  },
  data() {
    return {
      employees: [],
      loading: false,
      showModal: false,
      showEquipmentModal: false,
      modalMode: 'create',
      currentEmployee: this.getEmptyEmployee(),
      selectedEmployee: null,
      searchQuery: '',
      notification: {
        show: false,
        message: '',
        type: 'success'
      }
    }
  },
  methods: {
    getEmptyEmployee() {
      return {
        id: 0,
        name: '',
        position: '',
        department: '',
        contact: ''
      }
    },

    async loadData() {
      this.loading = true
      try {
        const response = await GetAllEmployees()
        if (response.model) {
          this.employees = response.model
        }
      } catch (error) {
        console.error('Ошибка загрузки сотрудников:', error)
        this.showNotification('Ошибка загрузки сотрудников', 'error')
      } finally {
        this.loading = false
      }
    },

    openCreateModal() {
      this.modalMode = 'create'
      this.currentEmployee = this.getEmptyEmployee()
      this.showModal = true
    },

    editEmployee(employee) {
      this.modalMode = 'edit'
      this.currentEmployee = { ...employee }
      this.showModal = true
    },

    viewEquipment(employee) {
      this.selectedEmployee = employee
      this.showEquipmentModal = true
    },

    async deleteEmployee(id) {
      if (!confirm('Вы уверены, что хотите удалить этого сотрудника?')) return

      try {
        const response = await DeleteEmployee(id)
        if (response.model) {
          this.showNotification('Сотрудник успешно удален')
          await this.loadData()
        }
      } catch (error) {
        console.error('Ошибка удаления сотрудника:', error)
        this.showNotification('Ошибка удаления сотрудника', 'error')
      }
    },

    async saveEmployee() {
      try {
        const service = this.modalMode === 'create' ? CreateEmployee : UpdateEmployee
        const response = await service(this.currentEmployee)
        
        if (response.model) {
          this.showNotification(
            this.modalMode === 'create'
              ? 'Сотрудник успешно добавлен'
              : 'Данные сотрудника обновлены'
          )
          this.closeModal()
          await this.loadData()
        }
      } catch (error) {
        console.error('Ошибка сохранения сотрудника:', error)
        this.showNotification('Ошибка сохранения сотрудника', 'error')
      }
    },

    closeModal() {
      this.showModal = false
      this.currentEmployee = this.getEmptyEmployee()
    },

    closeEquipmentModal() {
      this.showEquipmentModal = false
      this.selectedEmployee = null
    },

    applyFilters() {
      // Implement search functionality if needed
    },

    refreshData() {
      this.loadData()
    },

    showNotification(message, type = 'success') {
      this.notification = {
        show: true,
        message,
        type
      }
      setTimeout(() => {
        this.notification.show = false
      }, 3000)
    }
  },
  mounted() {
    this.loadData()
  }
}
</script>

<style scoped>
.employees-container {
  padding: 24px;
  background: #f8fafc;
  min-height: 100vh;
}

.header {
  margin-bottom: 32px;
}

.title {
  font-size: 32px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.action-panel {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
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

.table-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.employees-table {
  width: 100%;
  border-collapse: collapse;
}

.employees-table th {
  background: #f8fafc;
  padding: 16px 12px;
  text-align: left;
  font-weight: 600;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
  font-size: 14px;
}

.employees-table td {
  padding: 16px 12px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  font-size: 14px;
}

.actions {
  display: flex;
  gap: 8px;
}

.btn-icon {
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

.btn-icon:hover {
  background: #e2e8f0;
  color: #475569;
}

.btn-icon svg {
  width: 16px;
  height: 16px;
  stroke-width: 2;
}

.btn-link {
  background: none;
  border: none;
  color: #3b82f6;
  cursor: pointer;
  font-size: 14px;
  padding: 0;
  text-decoration: underline;
}

.btn-link:hover {
  color: #2563eb;
}

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

.modal-body {
  padding: 24px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.equipment-table {
  width: 100%;
  border-collapse: collapse;
}

.equipment-table th {
  background: #f8fafc;
  padding: 12px;
  text-align: left;
  font-weight: 600;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
  font-size: 14px;
}

.equipment-table td {
  padding: 12px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  font-size: 14px;
}

.notification {
  position: fixed;
  bottom: 24px;
  right: 24px;
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 14px;
  z-index: 1000;
  animation: slideIn 0.3s ease-out;
}

.notification-success {
  background: #dcfce7;
  color: #166534;
}

.notification-error {
  background: #fee2e2;
  color: #991b1b;
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

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px;
  color: #64748b;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e2e8f0;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style> 
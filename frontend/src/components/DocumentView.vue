<template>
  <div class="document-view">
    <v-card>
      <v-card-title class="d-flex justify-space-between align-center">
        <span>{{ title }}</span>
        <div>
          <v-btn
            v-if="document && document.status === 'draft'"
            color="primary"
            @click="$emit('edit')"
            class="mr-2"
          >
            <v-icon left>mdi-pencil</v-icon>
            Редактировать
          </v-btn>
          <v-btn
            v-if="document"
            color="secondary"
            @click="showExport = true"
          >
            <v-icon left>mdi-file-export</v-icon>
            Экспорт
          </v-btn>
        </div>
      </v-card-title>

      <!-- ... rest of the template ... -->
    </v-card>

    <!-- Dialog for export -->
    <v-dialog v-model="showExport" max-width="1200">
      <v-card>
        <v-card-title class="d-flex justify-space-between align-center">
          <span>Экспорт документа</span>
          <v-btn icon @click="showExport = false">
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>
        <v-card-text>
          <document-export
            v-if="showExport && document"
            :document-id="document.id"
          />
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import DocumentExport from './DocumentExport.vue'

export default {
  name: 'DocumentView',
  components: {
    DocumentExport
  },
  props: {
    document: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      showExport: false
    }
  },
  computed: {
    title() {
      if (!this.document) return 'Документ'
      const type = {
        'inventory': 'Инвентаризация',
        'transfer': 'Перемещение',
        'write_off': 'Списание',
        'acceptance': 'Приемка'
      }[this.document.type] || 'Документ'
      return `${type} №${this.document.number}`
    }
  }
}
</script>

<style scoped>
.document-view {
  width: 100%;
}
</style> 
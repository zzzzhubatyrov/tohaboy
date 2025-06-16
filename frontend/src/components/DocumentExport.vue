<template>
  <div class="document-export">
    <div v-if="loading" class="loading">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
      <span>Загрузка документа...</span>
    </div>
    <div v-else-if="error" class="error">
      <v-alert type="error">{{ error }}</v-alert>
    </div>
    <div v-else>
      <div class="actions mb-4">
        <v-btn color="primary" @click="print">
          <v-icon left>mdi-printer</v-icon>
          Печать
        </v-btn>
        <v-btn color="secondary" class="ml-2" @click="downloadTxt">
          <v-icon left>mdi-download</v-icon>
          Скачать TXT
        </v-btn>
      </div>
      <div class="document-content" ref="content">
        <pre>{{ content }}</pre>
      </div>
    </div>
  </div>
</template>

<script>
import { ExportDocument } from '../../wailsjs/go/service/DocumentService'

export default {
  name: 'DocumentExport',
  props: {
    documentId: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      content: '',
      loading: true,
      error: null
    }
  },
  async created() {
    try {
      this.content = await ExportDocument(this.documentId)
      this.loading = false
    } catch (error) {
      this.error = error.message || 'Ошибка при загрузке документа'
      this.loading = false
    }
  },
  methods: {
    print() {
      const printContent = this.content
      const printWindow = window.open('', '_blank')
      printWindow.document.write(`
        <html>
          <head>
            <title>Печать документа</title>
            <style>
              body { font-family: monospace; white-space: pre; }
            </style>
          </head>
          <body>
            ${printContent}
          </body>
        </html>
      `)
      printWindow.document.close()
      printWindow.print()
    },
    downloadTxt() {
      const blob = new Blob([this.content], { type: 'text/plain' })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = 'document.txt'
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
    }
  }
}
</script>

<style scoped>
.document-export {
  padding: 20px;
}

.loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  min-height: 200px;
}

.document-content {
  background: white;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.document-content pre {
  white-space: pre;
  font-family: monospace;
  margin: 0;
  overflow-x: auto;
}

.actions {
  display: flex;
  gap: 8px;
}
</style> 
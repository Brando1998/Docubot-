<template>
  <div class="max-w-4xl mx-auto p-6">
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-2">WhatsApp Bot</h1>
      <p class="text-gray-600">Administra la conexión de WhatsApp para el chatbot</p>
    </div>

    <!-- Estado de conexión -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-6">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-semibold text-gray-800">Estado de Conexión</h2>
        <div :class="statusBadgeClasses">
          <div :class="statusDotClasses" class="w-2 h-2 rounded-full mr-2"></div>
          {{ statusText }}
        </div>
      </div>

      <!-- WhatsApp conectado -->
      <div v-if="whatsappData?.connected" class="bg-green-50 border border-green-200 rounded-lg p-6 text-center">
        <div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-green-600" fill="currentColor" viewBox="0 0 24 24">
            <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.890-5.335 11.893-11.893A11.821 11.821 0 0020.885 3.488"/>
          </svg>
        </div>
        <h4 class="text-lg font-medium text-green-800 mb-2">WhatsApp Conectado</h4>
        <p class="text-green-600 text-sm mb-4">
          <strong>{{ whatsappData.user_name || 'Usuario' }}</strong> 
          ({{ whatsappData.phone_number || 'Sin número' }})
        </p>
        
        <div class="text-sm text-green-600 mb-4">
          <p><strong>Última conexión:</strong> {{ formatDate(whatsappData.last_connected) }}</p>
        </div>
        
        <button
          @click="handleDisconnect"
          :disabled="isLoading"
          class="bg-red-500 text-white px-6 py-2 rounded-lg hover:bg-red-600 disabled:opacity-50 flex items-center mx-auto"
        >
          <svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ isLoading ? 'Desconectando...' : 'Finalizar Sesión' }}
        </button>
      </div>
      
      <!-- QR Code para nueva conexión -->
      <div v-else-if="whatsappData?.status === 'waiting_for_scan'" class="bg-blue-50 border border-blue-200 rounded-lg p-6 text-center">
        <div class="mb-4">
          <h4 class="text-lg font-medium text-blue-800 mb-2">Escanea el Código QR</h4>
          <p class="text-blue-600 text-sm">Abre WhatsApp en tu teléfono y escanea este código</p>
        </div>
        
        <div class="bg-white p-4 rounded-lg border-2 border-dashed border-blue-300 inline-block mb-4">
          <div v-if="whatsappData.qr_image" class="w-64 h-64 flex items-center justify-center">
            <img 
              :src="whatsappData.qr_image" 
              alt="WhatsApp QR Code" 
              class="max-w-full max-h-full border border-gray-200 rounded"
            >
          </div>
          <div v-else class="w-64 h-64 flex items-center justify-center bg-gray-100 rounded">
            <div class="text-center">
              <svg class="w-8 h-8 text-gray-400 mx-auto mb-2" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clip-rule="evenodd"/>
              </svg>
              <p class="text-gray-500 text-sm">Generando QR...</p>
            </div>
          </div>
        </div>
        
        <button
          @click="handleRefreshQR"
          :disabled="isLoading"
          class="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600 disabled:opacity-50 flex items-center mx-auto"
        >
          <svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ isLoading ? 'Actualizando...' : 'Actualizar QR' }}
        </button>
      </div>

      <!-- Estado inicializando -->
      <div v-else-if="whatsappData?.status === 'initializing'" class="bg-yellow-50 border border-yellow-200 rounded-lg p-6 text-center">
        <div class="w-16 h-16 bg-yellow-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-yellow-600 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>
        <h4 class="text-lg font-medium text-yellow-800 mb-2">Inicializando WhatsApp</h4>
        <p class="text-yellow-600 text-sm">Configurando conexión, espera un momento...</p>
      </div>
      
      <!-- Estado desconectado -->
      <div v-else class="bg-gray-50 border border-gray-200 rounded-lg p-6 text-center">
        <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-gray-400" fill="currentColor" viewBox="0 0 24 24">
            <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.890-5.335 11.893-11.893A11.821 11.821 0 0020.885 3.488"/>
          </svg>
        </div>
        <h4 class="text-lg font-medium text-gray-800 mb-2">WhatsApp Desconectado</h4>
        <p class="text-gray-600 text-sm mb-4">Genera un código QR para conectar WhatsApp</p>
        
        <button
          @click="handleGenerateQR"
          :disabled="isLoading"
          class="bg-green-500 text-white px-6 py-2 rounded-lg hover:bg-green-600 disabled:opacity-50 flex items-center mx-auto"
        >
          <svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ isLoading ? 'Generando...' : 'Generar Código QR' }}
        </button>
      </div>
    </div>
    
    <!-- Mensajes de error/éxito -->
    <div v-if="message" :class="messageClasses" class="rounded-lg p-4 mb-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg v-if="messageType === 'success'" class="h-5 w-5 text-green-400" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
          </svg>
          <svg v-else class="h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
          </svg>
        </div>
        <div class="ml-3">
          <p class="text-sm font-medium">{{ message }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useWhatsApp } from '@/composables/useWhatsApp'

// Usar el composable
const {
  whatsappData,
  isLoading,
  error,
  fetchWhatsAppStatus,
  generateQR,
  refreshQR,
  disconnectWhatsApp,
  clearError
} = useWhatsApp()

// Estado para mensajes de UI
const message = ref('')
const messageType = ref<'success' | 'error'>('success')

// Estados computados
const statusText = computed(() => {
  if (!whatsappData.value) return 'Verificando...'
  if (whatsappData.value.connected) return 'Conectado'
  if (whatsappData.value.status === 'waiting_for_scan') return 'Esperando escaneo'
  if (whatsappData.value.status === 'initializing') return 'Inicializando'
  return 'Desconectado'
})

const statusBadgeClasses = computed(() => {
  const base = 'flex items-center px-3 py-1 rounded-full text-sm font-medium'
  if (!whatsappData.value) return `${base} bg-gray-100 text-gray-800`
  if (whatsappData.value.connected) return `${base} bg-green-100 text-green-800`
  if (whatsappData.value.status === 'waiting_for_scan') return `${base} bg-blue-100 text-blue-800`
  if (whatsappData.value.status === 'initializing') return `${base} bg-yellow-100 text-yellow-800`
  return `${base} bg-red-100 text-red-800`
})

const statusDotClasses = computed(() => {
  if (!whatsappData.value) return 'bg-gray-400'
  if (whatsappData.value.connected) return 'bg-green-400'
  if (whatsappData.value.status === 'waiting_for_scan') return 'bg-blue-400'
  if (whatsappData.value.status === 'initializing') return 'bg-yellow-400'
  return 'bg-red-400'
})

const messageClasses = computed(() => {
  const base = 'border'
  if (messageType.value === 'success') return `${base} bg-green-50 border-green-200`
  return `${base} bg-red-50 border-red-200`
})

// Handlers que usan el composable
const handleGenerateQR = async () => {
  try {
    await generateQR()
    showMessage('QR generado correctamente', 'success')
  } catch (err: any) {
    showMessage(error.value || 'Error generando QR', 'error')
  }
}

const handleRefreshQR = async () => {
  try {
    await refreshQR()
    showMessage('QR actualizado', 'success')
  } catch (err: any) {
    showMessage(error.value || 'Error actualizando QR', 'error')
  }
}

const handleDisconnect = async () => {
  try {
    await disconnectWhatsApp()
    showMessage('WhatsApp desconectado correctamente', 'success')
  } catch (err: any) {
    showMessage(error.value || 'Error desconectando', 'error')
  }
}

// Utilidades
const showMessage = (text: string, type: 'success' | 'error' = 'success') => {
  message.value = text
  messageType.value = type
  setTimeout(clearMessage, 5000)
}

const clearMessage = () => {
  message.value = ''
  clearError()
}

const formatDate = (dateString?: string) => {
  if (!dateString || dateString === '0001-01-01T00:00:00Z') return 'No disponible'
  return new Date(dateString).toLocaleString()
}

// Auto-refresh del estado cada 5 segundos
let refreshInterval: number | null = null

const startAutoRefresh = () => {
  refreshInterval = setInterval(() => {
    if (!isLoading.value) {
      fetchWhatsAppStatus()
    }
  }, 5000)
}

const stopAutoRefresh = () => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
    refreshInterval = null
  }
}

// Lifecycle
onMounted(() => {
  fetchWhatsAppStatus()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>
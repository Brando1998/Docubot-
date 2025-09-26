# 🤖 Docubot - Arquitectura y Flujo de Integración

## 📋 Descripción General

Docubot es un ecosistema de chatbot modular que combina múltiples tecnologías para automatizar procesos documentales en el sector transporte, principalmente a través de WhatsApp.

## 🛠️ Prerrequisitos del Sistema

- Docker 24.0+ y Docker Compose V2
- Make (Linux/macOS) o WSL2 (Windows)
- Al menos 8GB RAM disponible
- Puertos libres: 3000-3002, 5005, 5432, 8080, 27017

## 🏗️ Arquitectura de Servicios

### Componentes Principales

| Servicio | Tecnología | Puerto | Descripción |
|----------|------------|--------|-------------|
| **Vue Dashboard** | Vue 3 + TypeScript + Vite | 3002 | Frontend administrativo para gestión y configuración |
| **API Backend** | Go + Gin | 8080 | Hub central de comunicación y lógica de negocio |
| **Baileys Gateway** | Node.js + Baileys | 3000 | Conexión directa con WhatsApp Web |
| **Rasa Bot** | Python + Rasa | 5005 | Motor de NLP y gestión de conversaciones |
| **Playwright Actions** | Node.js + Playwright | 3001 | Automatización web para acciones específicas |
| **PostgreSQL** | PostgreSQL | 5432 | Base de datos principal |
| **MongoDB** | MongoDB | 27017 | Base de datos para documentos y logs |

## 🔄 Flujo de Comunicación Detallado

### 1. Flujo de Configuración (Dashboard Web)

```plaintext
[Administrador] 
    ↓ (HTTP/REST)
[Vue Dashboard :3002]
    ↓ (API calls)
[API Go :8080]
    ↓ (queries)
[PostgreSQL/MongoDB]
```

**Proceso:**
1. El administrador accede al dashboard Vue
2. Puede escanear el código QR para vincular WhatsApp
3. Configura respuestas del bot, usuarios, etc.
4. Toda la configuración se almacena en las bases de datos

### 2. Flujo Principal de Mensajería (WhatsApp ↔ Chatbot)

```plaintext
[Usuario WhatsApp]
    ↓ (mensaje texto/multimedia)
[Baileys :3000] ──────────── WebSocket ──────────── [API Go :8080]
    ↑ (respuesta)                                        ↓ (análisis)
                                                    [Rasa :5005]
                                                         ↓ (acciones)
                                                    [Playwright :3001]
                                                         ↓ (resultados)
                                                    [Websites Externos]
```

**Proceso paso a paso:**

#### A. Recepción de Mensaje
1. Usuario envía mensaje por WhatsApp
2. Baileys recibe el mensaje y lo procesa
3. Baileys envía via **WebSocket** a la API Go:
   ```json
   {
     "phone": "573001234567@s.whatsapp.net",
     "message": "Necesito un manifiesto",
     "botNumber": "573009876543@s.whatsapp.net"
   }
   ```

#### B. Procesamiento Central
4. API Go recibe el mensaje via WebSocket
5. API procesa y guarda en base de datos (cliente, bot, mensaje)
6. API envía mensaje a Rasa para análisis NLP:
   ```json
   {
     "sender": "573001234567",
     "message": "Necesito un manifiesto"
   }
   ```

#### C. Análisis y Respuesta
7. Rasa analiza el mensaje y determina:
   - **Intent**: solicitar_manifiesto
   - **Entities**: tipo de documento
   - **Action**: si requiere ejecutar acción especial
8. Rasa puede:
   - Devolver respuesta directa, O
   - Ejecutar acción personalizada con Playwright
9. Si requiere acción, Rasa llama a Playwright para automatizar websites
10. Rasa devuelve respuesta a la API:
    ```json
    [
      {
        "recipient_id": "573001234567",
        "text": "Perfecto, te ayudo con el manifiesto. ¿Para qué ruta necesitas el documento?"
      }
    ]
    ```

#### D. Envío de Respuesta
11. API Go procesa la respuesta de Rasa
12. API guarda la respuesta en base de datos
13. API envía via **WebSocket** a Baileys:
    ```json
    {
      "to": "573001234567@s.whatsapp.net",
      "message": "Perfecto, te ayudo con el manifiesto..."
    }
    ```
14. Baileys envía el mensaje de vuelta al usuario por WhatsApp

## 🔌 Detalles de Integración

### Comunicación WebSocket (Baileys ↔ API)

**Baileys → API:**
```typescript
// baileys-ws/src/handlers/messageHandler.ts
backendWS.send(JSON.stringify({
    phone: from,
    message: text,
    botNumber
}));
```

**API → Baileys:**
```go
// api/controllers/conversation.go
hub.SendToBot(msg.BotNumber, map[string]interface{}{
    "to":      msg.Phone,
    "message": response.Text,
})
```

### Comunicación HTTP (API ↔ Rasa)

**API → Rasa:**
```go
// POST http://rasa:5005/webhooks/rest/webhook
{
    "sender": "user_id",
    "message": "texto del mensaje"
}
```

**Rasa → API:**
```json
[
    {
        "recipient_id": "user_id", 
        "text": "respuesta del bot"
    }
]
```

### Acciones Personalizadas (Rasa ↔ Playwright)

Rasa puede ejecutar acciones personalizadas definidas en `actions/actions.py`:

```python
# rasa-bot/actions/actions.py
class ActionExpedirManifiesto(Action):
    def run(self, dispatcher, tracker, domain):
        # Llama a Playwright para automatizar proceso
        result = call_playwright_action("expedir_manifiesto", datos)
        dispatcher.utter_message(text=f"Manifiesto generado: {result}")
        return []
```

## 📊 Flujo de Datos

### Base de Datos
- **PostgreSQL**: Usuarios, bots, mensajes, configuraciones
- **MongoDB**: Documentos generados, logs, archivos multimedia

### Estados de Sesión
- **Baileys**: Mantiene sesión activa de WhatsApp
- **Rasa**: Mantiene contexto de conversación por usuario
- **API**: Gestiona estados de todas las sesiones

## 🚀 Casos de Uso Principales

### 1. Escaneo de QR (Configuración Inicial)
```plaintext
[Admin] → [Vue Dashboard] → [API] → [Baileys] → [WhatsApp Web]
```

### 2. Consulta Simple
```plaintext
[Usuario] → [WhatsApp] → [Baileys] → [API] → [Rasa] → respuesta directa
```

### 3. Generación de Documento
```plaintext
[Usuario] → [WhatsApp] → [Baileys] → [API] → [Rasa] → [Playwright] → [Website] → documento generado
```

## ⚙️ Configuración de Servicios

### Variables de Entorno Críticas
```bash
# Comunicación entre servicios
RASA_URL=http://rasa:5005
PLAYWRIGHT_URL=http://playwright:3001
API_URL=http://api:8080
BAILEYS_PORT=3000

# Base de datos
POSTGRES_HOST=postgres
MONGO_URI=mongodb://mongodb:27017
```

### Endpoints Principales
- **Vue Dashboard**: http://localhost:3002
- **API Health**: http://localhost:8080/health
- **API Swagger**: http://localhost:8080/swagger/index.html
- **Rasa Status**: http://localhost:5005/status
- **Baileys Health**: http://localhost:3000/health

## 🔧 Comandos de Desarrollo

### Levantar Entorno Completo
```bash
make up-sequential  # Recomendado: despliega servicios en orden
make up-local      # Alternativo: despliega todo junto
```

### Logs y Debugging
```bash
make logs-all       # Todos los logs
make logs-baileys   # Solo Baileys
make logs-api       # Solo API
make logs-rasa      # Solo Rasa
```

### Reiniciar Servicios
```bash
make restart-api
make restart-baileys  
make restart-rasa
```

## 🚨 Puntos Críticos

### Dependencias de Inicio
1. **PostgreSQL/MongoDB** deben estar listos primero
2. **API** debe iniciarse antes que Baileys
3. **Rasa** debe estar entrenado con modelo actual
4. **Baileys** necesita sesión activa de WhatsApp

### Gestión de Errores
- API maneja reconexiones automáticas con Rasa
- Baileys se reconecta automáticamente a WhatsApp
- Timeouts configurables para todas las comunicaciones
- Logs centralizados para debugging

### Seguridad
- Autenticación vía PASETO tokens
- Validación de mensajes entrantes
- Rate limiting por usuario
- Sanitización de datos antes de enviar a servicios externos

## 📈 Escalabilidad

### Horizontal
- Múltiples instancias de API detrás de load balancer
- Instancias separadas de Rasa por modelo/idioma
- Clustering de MongoDB para alta disponibilidad

### Vertical  
- Optimización de memoria en contenedores Go
- Caché Redis para sesiones frecuentes
- Optimización de queries en PostgreSQL

---

*Esta documentación refleja la implementación actual del sistema basada en el análisis del código fuente.*
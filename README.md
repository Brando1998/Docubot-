# 📦 Docubot

Docubot es un ecosistema modular para automatizar procesos documentales en el sector transporte, combinando:

- Backend en **Go**
- Dashboard frontend en **Vue.js** con TypeScript
- Bots conversacionales en **Rasa**
- Automatización web con **Playwright**
- Comunicación por WhatsApp con **Baileys**

---

## 🚀 Características

- 🔧 **API** RESTful en Go con PostgreSQL y MongoDB  
- 🎨 **Dashboard Vue.js** con TypeScript y Vite para gestión administrativa
- 🤖 **Chatbot** en Rasa con respuestas dinámicas y acciones personalizadas  
- 🧠 **Automatización web** en Node.js usando Playwright  
- 💬 **Integración WhatsApp** usando Baileys (no oficial)  
- ☸️ **Despliegue en Kubernetes** y local con Docker Compose  

---

## 📁 Estructura del Proyecto

```plaintext
.
├── api/
│   ├── cmd/api/
│   ├── config/
│   ├── controllers/
│   ├── databases/
│   ├── docs/
│   ├── middleware/
│   ├── mocks/
│   ├── models/
│   ├── repositories/
│   ├── routes/
│   └── services/
│
├── vue-dashboard/          # 🎨 Dashboard Vue.js con TypeScript
│   ├── src/
│   │   ├── components/
│   │   ├── assets/
│   │   └── main.ts
│   ├── public/
│   ├── package.json
│   ├── vite.config.ts
│   └── tsconfig*.json
│
├── baileys-ws/
│   ├── auth/               # ⚠️ Contiene credenciales/sesiones (ignorar en git)
│   ├── src/
│   │   ├── handlers/
│   │   ├── sessions/       # ⚠️ Contiene session-*.json (ignorar en git)
│   │   └── websocket/
│   └── package.json
│
├── docker/
│   ├── api.Dockerfile
│   ├── vue.Dockerfile      # 🎨 Dockerfile para Vue dashboard
│   └── baileys.Dockerfile
│
├── k8s/
│   ├── configmaps/
│   ├── deployments/
│   ├── secrets/
│   └── services/
│
├── playwright-bot/
│
├── rasa-bot/
│   ├── actions/
│   ├── data/
│   ├── models/             # ⚠️ Contiene modelos .tar.gz (ignorar en git)
│   └── requirements.txt
│
├── docker-compose.yml
├── Makefile
├── README.md
└── structure.md
```

---

## 🖥️ Requisitos

- Docker + Docker Compose  
- Make (opcional)  
- Kubernetes (Minikube o clúster compatible)  
- `kubectl` configurado  
- Git  
- Node.js 20+ (para desarrollo local del dashboard Vue)

---

## 📥 Clonar el Repositorio

```bash
git clone https://github.com/tuusuario/docubot.git
cd docubot
```

---

## ▶️ Ejecutar en Local (Docker Compose)

### 1. Configura variables de entorno
```bash
# Copia las variables de entorno
cp .env.example .env

# Edita .env con tus valores específicos
nano .env
```

### 2. Levanta el entorno local
```bash
# Opción 1: Despliegue secuencial (recomendado)
make up-sequential

# Opción 2: Despliegue completo 
make up-local
```

Esto iniciará:

- **Vue Dashboard** → http://localhost:3002  
- **API Go** → http://localhost:8080  
- **Rasa** → http://localhost:5005  
- **Playwright** (Node.js) → http://localhost:3001  
- **Baileys WS** (WhatsApp) → http://localhost:3000  
- **PostgreSQL** → localhost:5432  
- **MongoDB** → localhost:27017  

---

## 🌐 URLs de Acceso

| Servicio | URL | Descripción |
|----------|-----|-------------|
| **Dashboard Vue** | http://localhost:3002 | Interface administrativa |
| **API Go** | http://localhost:8080/health | Backend principal |
| **Rasa** | http://localhost:5005/status | Motor de chatbot |
| **Playwright** | http://localhost:3001/health | Automatización web |
| **Baileys** | http://localhost:3000/health | Gateway WhatsApp |

---

## ☸️ Verificar Estado

```bash
# Ver estado de contenedores
make status

# Ver logs en tiempo real
make logs-all

# Verificar salud de servicios
make health-check
```

---

## 🛠️ Comandos Útiles (Makefile)

### Contenedores
```bash
make up-local         # Levanta entorno local con docker-compose
make up-sequential    # Despliegue secuencial (recomendado)
make down-local       # Detiene el entorno local
make build-all        # Construir todas las imágenes del proyecto
```

### Logs
```bash
make logs-api         # Logs del backend API
make logs-vue         # Logs del dashboard Vue
make logs-rasa        # Logs de Rasa
make logs-playwright  # Logs de Playwright
make logs-baileys     # Logs de Baileys
make logs-all         # Logs de todos los servicios
```

### Desarrollo
```bash
make dev-logs         # Logs en tiempo real de todos los servicios
make dev-shell-api    # Shell en contenedor API
make dev-shell-vue    # Shell en contenedor Vue
make restart-vue      # Reiniciar solo Vue
make restart-api      # Reiniciar solo API
```

---

## 🎨 Desarrollo del Dashboard Vue

### Desarrollo local (fuera de Docker)
```bash
cd vue-dashboard

# Instalar dependencias
npm install

# Servidor de desarrollo
npm run dev

# Compilar para producción
npm run build
```

### Desarrollo con Docker
```bash
# Build y ejecutar contenedor Vue
make build-vue
make up-vue

# Ver logs de Vue
make logs-vue
```

---

## 🧪 Pruebas y Endpoints

### API Endpoints
- **Health Check**: http://localhost:8080/health  
- **Documentation**: http://localhost:8080/swagger/index.html  

### Chatbot
- **Rasa Webhook**: http://localhost:5005/webhooks/rest/webhook  
- **WhatsApp**: Se conecta automáticamente si está autenticado  

### Frontend
- **Dashboard**: http://localhost:3002  
- **Health Check**: http://localhost:3002/health  

---

## 🧩 Flujo de Integración

```plaintext
[Usuario Web] ↔ [Vue Dashboard] ↔ [API Go] ↔ [PostgreSQL/MongoDB]
                                      ↓
[Usuario WhatsApp] ↔ [Baileys-WS] ↔ [Rasa Bot] ↔ [Playwright Bot]
```

---

## 🏗️ Arquitectura de Servicios

### Frontend (Vue.js)
- **Tecnología**: Vue 3 + TypeScript + Vite
- **Puerto**: 3002
- **Servidor**: Nginx (producción)
- **Hot Reload**: Disponible en desarrollo

### Backend (Go)
- **Puerto**: 8080
- **Databases**: PostgreSQL + MongoDB
- **Documentation**: Swagger/OpenAPI

### Chatbot (Rasa)
- **Puerto**: 5005  
- **NLU**: Procesamiento de lenguaje natural
- **Actions**: Acciones personalizadas con Playwright

### WhatsApp Integration
- **Baileys**: Puerto 3000
- **WebSocket**: Comunicación en tiempo real

---

## 📚 Documentación API

```bash
# Generar documentación Swagger
cd api
swag init --output ./docs --dir ./cmd/api,./controllers
```

---

## ✅ Tests

```bash
# Tests del backend
cd api
go test ./controllers

# Tests del frontend (cuando estén implementados)
cd vue-dashboard
npm test
```

---

## 🐳 Docker

### Imágenes disponibles
- `docubot-api`: Backend en Go
- `docubot-vue`: Dashboard Vue.js
- `docubot-rasa`: Chatbot Rasa
- `docubot-playwright`: Automatización web
- `docubot-baileys`: Gateway WhatsApp

### Comandos Docker útiles
```bash
# Ver contenedores del proyecto
make show-containers

# Ver imágenes del proyecto  
make show-images

# Limpiar recursos del proyecto
make clean-project
```

---

## ☸️ Kubernetes

```bash
# Desplegar en Kubernetes
make k8s-deploy

# Eliminar despliegue
make k8s-delete
```

---

## 📝 Estructura de Desarrollo

### Agregar nuevas funcionalidades

#### Frontend (Vue)
1. Desarrollar en `vue-dashboard/src/`
2. Usar TypeScript para type safety
3. Compilar con `npm run build`
4. Testear con Docker: `make build-vue && make up-vue`

#### Backend (API)
1. Agregar endpoints en `api/controllers/`
2. Actualizar modelos en `api/models/`
3. Documentar con Swagger
4. Testear: `go test ./controllers`

#### Chatbot (Rasa)
1. Agregar intents en `rasa-bot/data/nlu.yml`
2. Definir stories en `rasa-bot/data/stories.yml`
3. Implementar acciones en `rasa-bot/actions/actions.py`
4. Reentrenar: `rasa train`

---

## 🔧 Configuración

### Variables de Entorno Principales

```bash
# Puertos
PORT=8080                    # API Backend
VUE_PORT=3002               # Dashboard Vue
RASA_PORT=5005              # Rasa Bot
PLAYWRIGHT_PORT=3001        # Playwright
BAILEYS_PORT=3000           # WhatsApp Gateway

# Base de datos
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DB=docubot_db
MONGO_URI=mongodb://mongodb:27017

# Servicios
RASA_URL=http://rasa:5005
PLAYWRIGHT_URL=http://playwright:3001
API_URL=http://api:8080
```

---

## 👨‍💻 Contribuciones

1. Haz fork del proyecto  
2. Crea una rama: `git checkout -b feature/nombre`  
3. Realiza tus cambios y haz commit: `git commit -m "Agrega nueva funcionalidad"`  
4. Push: `git push origin feature/nombre`  
5. Abre un Pull Request 🚀  

---

## 📄 Licencia

Este proyecto está bajo la licencia **MIT**.  

---

## 📬 Contacto

**Brando Díaz**  
✉️ brandodiazmont@gmail.com  
📱 WhatsApp +57 3023687930  
🔗 [LinkedIn](https://linkedin.com/in/brandodiaz)  

---

## 🔀 Flujo General Detallado

```plaintext
[Usuario Web Browser]
        ↓ (HTTP)
[Vue Dashboard :3002]
        ↓ (API calls)
[API Go :8080]
        ↓ (queries)
[PostgreSQL/MongoDB]

[Usuario WhatsApp]
        ↓ (mensaje)
[Baileys :3000]
        ↓ (webhook)
[API Go :8080]
        ↓ (NLU)
[RASA :5005]
        ↓ (actions)
[Playwright :3001] → [Websites]
```
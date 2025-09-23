# 📦 Docubot

Docubot es un ecosistema modular para automatizar procesos documentales en el sector transporte, combinando:

- Backend en **Go**
- Bots conversacionales en **Rasa**
- Automatización web con **Playwright**
- Comunicación por WhatsApp con **Baileys**

---

## 🚀 Características

- 🔧 **API** RESTful en Go con PostgreSQL y MongoDB  
- 🤖 **Chatbot** en Rasa con respuestas dinámicas y acciones personalizadas  
- 🧠 **Automatización web** en Node.js usando Playwright  
- 💬 **Integración WhatsApp** usando Baileys (no oficial)  
- ☸️ **Despliegue en Kubernetes** y local con Docker Compose  

---

## 📁 Estructura del Proyecto

```plaintext
.
├── api
│   ├── cmd/api
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
├── baileys-ws
│   ├── auth/           # ⚠️ Contiene credenciales/sesiones (ignorar en git)
│   ├── src/
│   │   ├── handlers/
│   │   ├── sessions/   # ⚠️ Contiene session-*.json (ignorar en git)
│   │   └── websocket/
│   └── package.json
│
├── docker/
│   └── api.Dockerfile
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
│   ├── models/         # ⚠️ Contiene modelos .tar.gz (ignorar en git)
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

- API Go → http://localhost:8000  
- Rasa → http://localhost:5005  
- Playwright (Node.js)  
- Baileys WS (WhatsApp)  
- PostgreSQL  
- MongoDB  

---

## ☸️ Verificar estado

```bash
# Ver estado de contenedores
make status

# Ver logs en tiempo real
make logs-all

# Verificar salud de servicios
make health-check
```

*(Opcional) Exponer vía Ingress con un controlador de entrada, o ajustar NodePort en los servicios.*  

---

## Urls de acceso

Rasa: http://localhost:5005/status
Playwright: http://localhost:3001/health
API: http://localhost:8000/health
Baileys: http://localhost:3000/health

## 🛠️ Comandos Útiles (Makefile)

```bash
make up-local       # Levanta entorno local con docker-compose
make down-local     # Detiene el entorno local
make k8s-deploy     # Aplica todos los manifiestos Kubernetes
make build-api      # Build del backend en Go
make logs-api       # Logs del contenedor API
```

---

## 🧪 Pruebas y Endpoints

- API: http://localhost:8000/health  
- Rasa: http://localhost:5005/webhooks/rest/webhook  
- WhatsApp (Baileys): se conecta automáticamente y responde si está autenticado  

---

## 🧩 Flujo de Integración

```plaintext
[Usuario] ↔ [Baileys-WS] ↔ [Rasa Bot] ↔ [API Go] ↔ [Playwright / DBs]
```

---

## 📚 Documentación API

```bash
swag init --output ./docs --dir ./cmd/api,./controllers
```

---

## ✅ Tests

```bash
go test ./controllers
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

Brando Díaz  
✉️ brandodiazmont@gmail.com  
📱 WhatsApp +57 3023687930  
🔗 LinkedIn  

---

## 🔀 Flujo General Detallado

```plaintext
[Usuario en WhatsApp]
        ↓ (mensaje)
[Baileys (Node.js)]
        ↓ (evento onMessage)
[API en Go]
        ↓ (POST /mensaje)
[RASA (NLU)]
        ↓ (intención/respuesta)
[API: evalúa acción]
    ┌─────────────┬─────────────┐
    |             |             |
[Guardar]   [Responder]   [Playwright]
en MongoDB  a WhatsApp     (Genera PDF, etc.)
                             ↓
                    [Documento generado]
                             ↓
                  [Baileys lo envía al usuario]
```

| Módulo         | Rol                                                             |
| -------------- | --------------------------------------------------------------- |
| **Baileys**    | Entrada y salida de mensajes con WhatsApp                       |
| **API (Go)**   | Centro de lógica, guarda todo, decide qué hacer                 |
| **MongoDB**    | Guarda mensajes, conversaciones, documentos                     |
| **Rasa**       | Interpreta qué quiere el usuario                                |
| **Playwright** | Genera documentos visuales que el bot puede entregar al cliente |
| **Postgres**   | Administra usuarios y datos estructurados (GORM)                |

---

## 📊 Métricas Sugeridas

- Último acceso de cada usuario  
- Cantidad de documentos generados  
- Cantidad de mensajes por bot  

*(Útil para administración o monetización del sistema)*  
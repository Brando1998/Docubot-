# 🎨 Vue Dashboard - Docubot

Dashboard administrativo desarrollado con **Vue.js 3**, **TypeScript** y **Vite** para la gestión y monitoreo del ecosistema Docubot.

---

## 🚀 Características

- ⚡ **Vue 3** con Composition API
- 🔷 **TypeScript** para type safety
- ⚡ **Vite** para desarrollo rápido y build optimizado
- 🎨 **Responsive Design** adaptable a todos los dispositivos
- 🔄 **Hot Module Replacement** para desarrollo ágil
- 📊 **Dashboard Interactivo** para monitoreo de servicios
- 🔗 **API Integration** con el backend Go

---

## 📁 Estructura del Proyecto

```plaintext
vue-dashboard/
├── public/                 # Archivos estáticos
│   └── vite.svg
├── src/
│   ├── components/         # Componentes Vue reutilizables
│   │   └── HelloWorld.vue
│   ├── assets/            # Recursos estáticos (imágenes, CSS)
│   │   └── vue.svg
│   ├── App.vue            # Componente raíz
│   ├── main.ts            # Punto de entrada de la aplicación
│   └── style.css          # Estilos globales
├── index.html             # HTML template
├── package.json           # Dependencias y scripts
├── package-lock.json      # Lock file de dependencias
├── tsconfig.app.json      # Config TypeScript para app
├── tsconfig.json          # Config TypeScript base
├── tsconfig.node.json     # Config TypeScript para Node
└── vite.config.ts         # Configuración de Vite
```

---

## 🛠️ Requisitos

- **Node.js** 20.19.0+ o >= 22.12.0
- **npm** o **yarn**
- **TypeScript** 5.8+

---

## ⚡ Desarrollo Local

### 1. Instalar dependencias
```bash
cd vue-dashboard
npm install
```

### 2. Servidor de desarrollo
```bash
# Inicia el servidor de desarrollo con HMR
npm run dev

# La aplicación estará disponible en:
# http://localhost:5173
```

### 3. Build para producción
```bash
# Compila la aplicación para producción
npm run build

# Los archivos compilados se guardan en ./dist/
```

### 4. Preview de producción
```bash
# Previsualiza la build de producción localmente
npm run preview
```

---

## 🐳 Desarrollo con Docker

### Build y ejecución
```bash
# Desde el directorio raíz del proyecto
make build-vue
make up-vue

# O específicamente:
docker compose up vue --build
```

### Ver logs
```bash
make logs-vue
```

### Abrir shell en contenedor
```bash
make dev-shell-vue
```

---

## 🔧 Configuración

### Variables de Entorno

La aplicación Vue puede configurarse mediante variables de entorno:

```bash
# Desarrollo
VITE_API_URL=http://localhost:8080
VITE_APP_NAME=Docubot Dashboard
VITE_APP_VERSION=1.0.0

# Producción (en Docker)
VUE_APP_API_URL=http://api:8080
NODE_ENV=production
```

### Configuración de Vite

El archivo `vite.config.ts` contiene la configuración del bundler:

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  // Configuraciones adicionales aquí
})
```

---

## 📦 Dependencias

### Principales
- **vue**: ^3.5.21 - Framework principal
- **typescript**: ~5.8.3 - Tipado estático
- **vite**: ^7.1.7 - Build tool y dev server

### DevDependencies
- **@vitejs/plugin-vue**: Soporte Vue para Vite
- **@vue/tsconfig**: Configuraciones TypeScript para Vue
- **vue-tsc**: Type-checking para Vue con TypeScript

---

## 🔗 Integración con API

El dashboard se conecta al backend Go a través de HTTP requests:

```typescript
// Ejemplo de llamada a la API
const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

async function fetchHealthStatus() {
  const response = await fetch(`${apiUrl}/health`)
  return response.json()
}
```

---

## 🎯 Funcionalidades Planificadas

### Dashboard Principal
- [ ] Monitoreo de servicios en tiempo real
- [ ] Métricas de performance del chatbot
- [ ] Logs centralizados
- [ ] Estado de conexiones WhatsApp

### Gestión de Manifiestos
- [ ] Listado de manifiestos generados
- [ ] Formulario de creación manual
- [ ] Búsqueda y filtrado avanzado
- [ ] Descarga de documentos PDF

### Configuración del Sistema
- [ ] Configuración de Rasa
- [ ] Gestión de respuestas del bot
- [ ] Configuración de Playwright
- [ ] Parámetros del sistema

### Analytics y Reportes
- [ ] Estadísticas de uso
- [ ] Reportes de conversaciones
- [ ] Métricas de automatización
- [ ] Dashboard ejecutivo

---

## 🎨 Personalización de UI

### Estilos
Los estilos están centralizados en:
- `src/style.css` - Estilos globales
- Componentes individuales con `<style scoped>`

### Temas
Para implementar temas personalizados:

```css
/* Variables CSS para theming */
:root {
  --primary-color: #646cff;
  --secondary-color: #42b883;
  --background-color: #242424;
  --text-color: rgba(255, 255, 255, 0.87);
}
```

---

## 🧪 Testing

```bash
# Tests (cuando estén implementados)
npm run test

# Tests con coverage
npm run test:coverage
```

---

## 📦 Build y Despliegue

### Build Local
```bash
npm run build
```

### Build con Docker
```bash
# Multi-stage build para producción
docker build -f ../docker/vue.Dockerfile -t docubot-vue .
```

### Nginx en Producción
El dashboard usa Nginx como servidor web en producción:
- Puerto: 3002
- Health check: `/health`
- SPA routing configurado
- Gzip compression habilitada
- Proxy API configurado en `/api/`

---

## 🔄 Scripts Disponibles

```json
{
  "dev": "vite",                    // Servidor de desarrollo
  "build": "vue-tsc -b && vite build", // Build para producción
  "preview": "vite preview"         // Preview de producción
}
```

---

## 🚨 Troubleshooting

### Errores Comunes

**Error: Module not found**
```bash
# Limpiar node_modules y reinstalar
rm -rf node_modules package-lock.json
npm install
```

**Error de tipos TypeScript**
```bash
# Verificar configuración TypeScript
npx vue-tsc --noEmit
```

**Error de puerto en uso**
```bash
# Verificar puertos ocupados
make check-ports
```

### Health Checks

```bash
# Verificar que Vue responde
curl http://localhost:3002/health

# En Docker
docker exec docubot-vue wget --spider http://localhost:3002/health
```

---

## 🔐 Consideraciones de Seguridad

- Headers de seguridad configurados en Nginx
- Variables de entorno para configuración
- Build de producción optimizado
- Sin exposición de API keys en frontend

---

## 🤝 Contribuir

### Estructura de Componentes
```
src/components/
├── common/          # Componentes comunes
├── dashboard/       # Componentes del dashboard
├── forms/           # Formularios
└── layout/          # Componentes de layout
```

### Convenciones
- Usar Composition API con `<script setup>`
- TypeScript obligatorio para type safety
- Props y emits tipados
- CSS scoped en componentes

---

## 📚 Recursos

- [Vue.js 3 Documentation](https://vuejs.org/)
- [TypeScript Vue Guide](https://vuejs.org/guide/typescript/overview.html)
- [Vite Documentation](https://vitejs.dev/)
- [Vue DevTools](https://devtools.vuejs.org/)

---

## 📄 Licencia

Este proyecto es parte de Docubot y está bajo la licencia **MIT**.

---

**Nota**: Este dashboard está en desarrollo activo. Las funcionalidades marcadas con ☑️ están implementadas, las marcadas con ◻️ están en desarrollo.
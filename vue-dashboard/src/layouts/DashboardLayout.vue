<template>
  <div class="flex h-screen">
    <!-- Sidebar -->
    <aside class="w-64 bg-gray-800 text-white flex flex-col">
      <div class="p-6 text-2xl font-bold border-b border-gray-700">
        Dashboard
      </div>
      <nav class="flex-1 p-4">
        <ul class="space-y-4">
          <li><router-link to="/dashboard" class="hover:text-blue-400">Inicio</router-link></li>
          <li><router-link to="/dashboard/profile" class="hover:text-blue-400">Perfil</router-link></li>
          <li><router-link to="/dashboard/settings" class="hover:text-blue-400">Configuración</router-link></li>
        </ul>
      </nav>
    </aside>

    <!-- Contenido principal -->
    <main class="flex-1 bg-gray-100 p-6 overflow-y-auto">
      <header class="flex justify-between items-center mb-6">
        <h1 class="text-3xl font-bold">Panel</h1>
        <button
          @click="handleLogout"
          class="bg-red-500 text-white px-4 py-2 rounded-lg hover:bg-red-600"
        >
          Cerrar sesión
        </button>
      </header>

      <!-- Renderiza la subvista -->
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from "@/composables/useAuth";
import { useRouter } from "vue-router";

const { logout } = useAuth();
const router = useRouter();

const handleLogout = () => {
  logout();
  router.push("/login");
};
</script>

<template>
  <div class="flex flex-col items-center justify-center h-screen text-center">
    <h1 class="text-6xl font-bold text-yellow-500">403</h1>
    <p class="text-xl mt-4">No tienes permisos para acceder a esta página</p>
    <RouterLink to="/login" class="mt-6 px-4 py-2 bg-blue-600 text-white rounded">
      Ir a Login
    </RouterLink>
  </div>
</template>

<script setup>
import { RouterLink } from "vue-router";
</script>

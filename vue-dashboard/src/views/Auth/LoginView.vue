<template>
  <div class="flex justify-center items-center h-screen bg-gray-100">
    <div class="bg-white shadow-lg rounded-2xl p-8 w-96">
      <h2 class="text-2xl font-bold mb-6 text-center">Login</h2>
      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label class="block mb-2 font-medium">Email</label>
          <input v-model="username" type="text" class="w-full p-2 border rounded-lg" required />
        </div>
        <div class="mb-6">
          <label class="block mb-2 font-medium">Password</label>
          <input v-model="password" type="password" class="w-full p-2 border rounded-lg" required />
        </div>
        <button
          type="submit"
          class="w-full bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600"
        >
          Iniciar Sesión
        </button>
      </form>
      <p v-if="errorMessage" class="text-red-500 mt-4">{{ errorMessage }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useAuth } from "@/composables/useAuth";
import { useRouter } from "vue-router";

const { login } = useAuth();
const router = useRouter();

const username = ref("");
const password = ref("");
const errorMessage = ref("");

const handleLogin = async () => {
  try {
    await login(username.value, password.value);
    router.push("/dashboard"); // ajusta tu ruta privada
  } catch (e) {
    errorMessage.value = "Credenciales inválidas";
  }
};
</script>

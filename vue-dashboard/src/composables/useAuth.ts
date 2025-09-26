import { ref, computed } from "vue";
import api from "../services/api";
import { decodePasetoPayload } from "../utils/token";

const accessToken = ref<string | null>(localStorage.getItem("accessToken"));
const refreshToken = ref<string | null>(localStorage.getItem("refreshToken"));

const payload = computed(() => {
  if (!accessToken.value) return null;
  return decodePasetoPayload(accessToken.value);
});

const userRole = computed(() => payload.value?.role || null);
const userEmail = computed(() => payload.value?.email || null);
const userId = computed(() => payload.value?.sub || null);

export function useAuth() {
  const login = async (username: string, password: string) => {
    const res = await api.post("/auth/login", { username, password });
    accessToken.value = res.data.access_token;
    refreshToken.value = res.data.refresh_token;

    localStorage.setItem("accessToken", accessToken.value!);
    localStorage.setItem("refreshToken", refreshToken.value!);
  };

  const logout = () => {
    accessToken.value = null;
    refreshToken.value = null;
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
  };

  const refreshTokenFn = async () => {
    if (!refreshToken.value) throw new Error("No refresh token");

    const res = await api.post("/auth/refresh", {
      refresh_token: refreshToken.value,
    });

    accessToken.value = res.data.access_token;
    localStorage.setItem("accessToken", accessToken.value!);
  };

  return {
    accessToken,
    refreshToken,
    payload,
    userRole,
    userEmail,
    userId,
    login,
    logout,
    refreshTokenFn,
  };
}

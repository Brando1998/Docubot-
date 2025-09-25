import axios from "axios";
import { useAuth } from "../composables/useAuth";

const api = axios.create({
  baseURL: "http://localhost:8080", // ajusta a tu backend
});

api.interceptors.request.use(
  (config) => {
    const { accessToken } = useAuth();
    if (accessToken.value) {
      config.headers.Authorization = `Bearer ${accessToken.value}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const { refreshTokenFn } = useAuth();

    if (error.response?.status === 401) {
      try {
        await refreshTokenFn();
        return api(error.config);
      } catch (e) {
        console.error("Refresh token failed", e);
      }
    }
    return Promise.reject(error);
  }
);

export default api;

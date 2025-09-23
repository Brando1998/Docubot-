// src/sessions/auth.ts
import { useMultiFileAuthState } from "@whiskeysockets/baileys";
/**
 * Obtiene el estado de autenticación usando múltiples archivos.
 * @returns {Promise<ReturnType<typeof useMultiFileAuthState>>} El estado de autenticación.
 */
export const getAuthState = async () => {
  const authState = await useMultiFileAuthState("auth");
  return authState;
};

export function decodePasetoPayload(token: string): Record<string, any> | null {
  try {
    const parts = token.split(".");
    if (parts.length < 3) return null; // no es un token válido

    const payloadBase64 = parts[2]; // 👈 Paseto guarda el payload aquí
    const padded = payloadBase64!.padEnd(payloadBase64!.length + (4 - (payloadBase64!.length % 4)) % 4, "=");
    const json = atob(padded.replace(/-/g, "+").replace(/_/g, "/"));

    return JSON.parse(json);
  } catch (err) {
    console.error("Error decoding Paseto:", err);
    return null;
  }
}

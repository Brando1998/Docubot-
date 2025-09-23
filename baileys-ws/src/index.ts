import { makeWASocket, fetchLatestBaileysVersion } from "@whiskeysockets/baileys";
import P from "pino";
import { getAuthState } from "./sessions/auth.js";
import qrcode from "qrcode-terminal";
import { handleIncomingMessage } from "./handlers/messageHandler.js";
import { connectToBackendWS } from "./websocket/client.js";

const start = async () => {
    // 1. Iniciar Baileys   
    const { state, saveCreds } = await getAuthState();
    const { version } = await fetchLatestBaileysVersion();

    const sock = makeWASocket({
        version,
        logger: P({ level: "silent" }),
        auth: state,
        defaultQueryTimeoutMs: undefined
    });

    // 2. Conectar al backend Go via WebSocket
    const backendWS = await connectToBackendWS(sock.user!.id); // Usa el número de teléfono del bot

    // 3. Eventos de Baileys
    sock.ev.on("connection.update", (update) => {
        const { connection, lastDisconnect, qr } = update;
        if (qr) qrcode.generate(qr, { small: true });
        if (connection === "open") console.log("✅ Conexión WhatsApp establecida");
    });

    sock.ev.on("creds.update", saveCreds);
    sock.ev.on("messages.upsert", async ({ messages }) => {
        for (const msg of messages) {
            if (!msg.message) continue;
            const from = msg.key.remoteJid;
            const text = msg.message.conversation || msg.message.extendedTextMessage?.text;
            const bot_number = sock.user!.id;

            if (!from || !text || msg.key.fromMe) continue;

            console.log('Mensaje recibido de', from, ':', text);

            try {
                await handleIncomingMessage(from, text, bot_number, backendWS);
                console.log('Mensaje enviado al backend:', {
                    phone: from,
                    message: text,
                    bot_number: bot_number
                });
            } catch (error) {
                console.error('Error enviando mensaje al backend:', error);
            }
        }
    });

    // 4. Escuchar respuestas del backend
    backendWS.on('message', (data) => {
        try {
            const response = JSON.parse(data.toString());
            console.log('Respuesta recibida del backend:', response);
            if (response.to && response.message) {
                sock.sendMessage(response.to, { text: response.message })
                    .then(() => console.log('Respuesta enviada a WhatsApp'))
                    .catch(err => console.error('Error enviando a WhatsApp:', err));
            }
        } catch (err) {
            console.error('Error procesando respuesta del backend:', err);
        }
    });
};

start().catch(err => {
    console.error("Error en la aplicación:", err);
    process.exit(1);
});
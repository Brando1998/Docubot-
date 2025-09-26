// baileys-ws/src/index.ts - Versión actualizada con endpoints HTTP
import "dotenv/config";
import { makeWASocket, fetchLatestBaileysVersion } from "@whiskeysockets/baileys";
import P from "pino";
import { getAuthState } from "./sessions/auth.js";
import qrcode from "qrcode-terminal";
import QRCode from "qrcode";
import { handleIncomingMessage } from "./handlers/messageHandler.js";
import { connectToBackendWS } from "./websocket/client.js";
import express from 'express';

const app = express();
const PORT = process.env.WS_PORT || 3000;

// Middleware
app.use(express.json());

// Variables globales para mantener estado
let sock: any = null;
let qrCodeData: string = '';
let currentStatus = {
    connected: false,
    number: '',
    name: '',
    qr_code: '',
    qr_image: ''
};

// =============================================
// ENDPOINTS HTTP PARA EL BACKEND
// =============================================

// Health check endpoint
app.get('/health', (req, res) => {
    res.json({ 
        status: 'ok', 
        service: 'baileys-ws',
        timestamp: new Date().toISOString(),
        whatsapp_connected: currentStatus.connected
    });
});

// Status endpoint
app.get('/status', (req, res) => {
    res.json({ 
        status: 'running',
        uptime: process.uptime(),
        memory: process.memoryUsage(),
        whatsapp: currentStatus
    });
});

// 🆕 Endpoint para obtener QR o estado
app.get('/qr', async (req, res) => {
    try {
        if (currentStatus.connected) {
            // Ya hay sesión activa
            res.json({
                status: 'connected',
                message: 'WhatsApp ya está conectado',
                connected: true,
                session_info: {
                    number: currentStatus.number,
                    name: currentStatus.name,
                    last_seen: new Date().toISOString()
                }
            });
        } else if (qrCodeData) {
            // Hay QR disponible
            const qrImageBase64 = await QRCode.toDataURL(qrCodeData);
            
            res.json({
                status: 'waiting_scan',
                message: 'Escanea el código QR con WhatsApp',
                connected: false,
                qr_code: qrCodeData,
                qr_image: qrImageBase64
            });
        } else {
            // Generar nuevo QR
            await restartBaileys();
            res.json({
                status: 'generating',
                message: 'Generando código QR...',
                connected: false
            });
        }
    } catch (error:any) {
        console.error('Error en /qr:', error);
        res.status(500).json({
            error: 'Error generando QR',
            details: error.message
        });
    }
});

// 🆕 Endpoint para desconectar
app.post('/disconnect', async (req, res) => {
    try {
        if (sock) {
            await sock.logout();
            sock.end();
            sock = null;
        }
        
        // Limpiar estado
        currentStatus = {
            connected: false,
            number: '',
            name: '',
            qr_code: '',
            qr_image: ''
        };
        qrCodeData = '';
        
        res.json({
            success: true,
            message: 'Sesión de WhatsApp terminada',
            status: 'disconnected'
        });
    } catch (error:any) {
        console.error('Error desconectando:', error);
        res.status(500).json({
            error: 'Error desconectando WhatsApp',
            details: error.message
        });
    }
});

// 🆕 Endpoint para enviar mensajes (desde el backend)
app.post('/send', async (req, res) => {
    try {
        const { to, message } = req.body;
        
        if (!currentStatus.connected || !sock) {
            return res.status(400).json({
                error: 'WhatsApp no está conectado'
            });
        }
        
        await sock.sendMessage(to, { text: message });
        
        res.json({
            success: true,
            message: 'Mensaje enviado correctamente',
            to: to
        });
    } catch (error:any) {
        console.error('Error enviando mensaje:', error);
        res.status(500).json({
            error: 'Error enviando mensaje',
            details: error.message
        });
    }
});

// 🆕 Endpoint para reiniciar/crear nueva sesión
app.post('/restart', async (req, res) => {
    try {
        console.log('🔄 Reiniciando sesión de WhatsApp...');
        await restartBaileys();
        
        res.json({
            success: true,
            message: 'Sesión reiniciada correctamente',
            status: 'restarting'
        });
    } catch (error:any) {
        console.error('Error reiniciando sesión:', error);
        res.status(500).json({
            error: 'Error reiniciando sesión',
            details: error.message
        });
    }
});

// =============================================
// FUNCIONES DE BAILEYS
// =============================================

const startHealthServer = () => {
    app.listen(PORT, () => {
        console.log(`🏥 Baileys HTTP Server running on port ${PORT}`);
    });
};

const restartBaileys = async () => {
    try {
        console.log('🔄 Reiniciando Baileys...');
        
        if (sock) {
            sock.end();
        }
        
        const { state, saveCreds } = await getAuthState();
        const { version } = await fetchLatestBaileysVersion();

        sock = makeWASocket({
            version,
            logger: P({ level: "silent" }),
            auth: state,
            defaultQueryTimeoutMs: undefined
        });

        setupBaileysEvents(sock, saveCreds);
        
    } catch (error) {
        console.error('Error reiniciando Baileys:', error);
    }
};

const setupBaileysEvents = (socket: any, saveCreds: any) => {
    // Evento de actualización de conexión
    socket.ev.on("connection.update", async (update: any) => {
        const { connection, lastDisconnect, qr } = update;
        
        if (qr) {
            qrCodeData = qr;
            qrcode.generate(qr, { small: true });
            console.log('📱 Nuevo QR generado');
            
            // Actualizar estado
            try {
                const qrImageBase64 = await QRCode.toDataURL(qr);
                currentStatus.qr_code = qr;
                currentStatus.qr_image = qrImageBase64;
            } catch (error) {
                console.error('Error generando QR image:', error);
            }
        }
        
        if (connection === "open") {
            console.log("✅ Conexión WhatsApp establecida");
            
            // Actualizar estado
            currentStatus.connected = true;
            currentStatus.number = socket.user?.id || '';
            currentStatus.name = socket.user?.name || 'Bot Docubot';
            qrCodeData = ''; // Limpiar QR
            
        } else if (connection === "close") {
            console.log("❌ Conexión WhatsApp cerrada");
            currentStatus.connected = false;
            
            // Intentar reconectar si no fue logout manual
            if (!lastDisconnect?.error?.message?.includes('logout')) {
                setTimeout(() => restartBaileys(), 3000);
            }
        }
    });

    // Guardar credenciales
    socket.ev.on("creds.update", saveCreds);
    
    // Manejar mensajes entrantes
    socket.ev.on("messages.upsert", async ({ messages }: any) => {
        for (const msg of messages) {
            if (!msg.message) continue;
            const from = msg.key.remoteJid;
            const text = msg.message.conversation || msg.message.extendedTextMessage?.text;
            const bot_number = socket.user!.id;

            if (!from || !text || msg.key.fromMe) continue;

            console.log('📨 Mensaje recibido de', from, ':', text);

            try {
                // Conectar al backend si no está conectado
                const backendWS = await connectToBackendWS(bot_number);
                await handleIncomingMessage(from, text, bot_number, backendWS);
                console.log('✅ Mensaje enviado al backend');
            } catch (error) {
                console.error('❌ Error enviando mensaje al backend:', error);
            }
        }
    });
};

// =============================================
// INICIALIZACIÓN
// =============================================

const start = async () => {
    try {
        // 1. Iniciar servidor HTTP
        startHealthServer();
        
        // 2. Iniciar Baileys
        await restartBaileys();
        
        console.log('🚀 Baileys-WS iniciado correctamente');
        
    } catch (error) {
        console.error("❌ Error iniciando aplicación:", error);
        process.exit(1);
    }
};

// Manejar cierre graceful
process.on('SIGINT', async () => {
    console.log('🛑 Cerrando Baileys...');
    if (sock) {
        await sock.logout();
        sock.end();
    }
    process.exit(0);
});

start();
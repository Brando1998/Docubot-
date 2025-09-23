import { WebSocket } from 'ws';

export const handleIncomingMessage = async (
    from: string,
    text: string,
    botNumber: string,
    backendWS: WebSocket
) => {
    // Enviar mensaje al backend Go
    backendWS.send(JSON.stringify({
        phone: from,
        message: text,
        botNumber
    }));

    console.log(`Mensaje enviado al backend: ${from} - ${text}`);
};
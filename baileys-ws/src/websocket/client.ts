import WebSocket from 'ws';

export const connectToBackendWS = (phone: string): Promise<WebSocket> => {
    return new Promise((resolve, reject) => {
        // ✅ Usar variable de entorno en lugar de localhost hardcodeado
        const apiUrl = process.env.API_URL || 'http://localhost:8080';
        const wsUrl = apiUrl.replace('http:', 'ws:').replace('https:', 'wss:');
        
        const ws = new WebSocket(`${wsUrl}/ws?phone=${phone}`);

        ws.on('open', () => {
            console.log('✅ Conectado al backend Go');
            resolve(ws);
        });

        ws.on('error', (err) => {
            console.error('Error conectando al backend:', err);
            reject(err);
        });

        ws.on('close', () => {
            console.log('Conexión WebSocket cerrada');
        });
    });
};
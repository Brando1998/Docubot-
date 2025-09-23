import WebSocket from 'ws';

export const connectToBackendWS = (phone: string): Promise<WebSocket> => {
    return new Promise((resolve, reject) => {
        const ws = new WebSocket(`ws://localhost:8080/ws?phone=${phone}`); // Añade el parámetro phone

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
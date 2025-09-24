import express from 'express';

const app = express();
const PORT = process.env.WS_PORT || 3000;

// Health check endpoint
app.get('/health', (req, res) => {
    res.json({ 
        status: 'ok', 
        service: 'baileys-ws',
        timestamp: new Date().toISOString()
    });
});

// Status endpoint
app.get('/status', (req, res) => {
    res.json({ 
        status: 'running',
        uptime: process.uptime(),
        memory: process.memoryUsage()
    });
});

export const startHealthServer = () => {
    app.listen(PORT, () => {
        console.log(`🏥 Health server running on port ${PORT}`);
    });
};
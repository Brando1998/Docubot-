const express = require('express');
const { chromium } = require('playwright');

const app = express();
const PORT = process.env.PORT || 3001;

app.use(express.json());

// Health check endpoint
app.get('/health', (req, res) => {
    res.json({ status: 'ok', service: 'playwright-bot' });
});

// Endpoint para generar manifiestos (ejemplo)
app.post('/generate-manifest', async (req, res) => {
    try {
        const { data } = req.body;
        
        console.log('Generando manifiesto con datos:', data);
        
        // Aquí irá la lógica de Playwright para automatizar la generación
        const browser = await chromium.launch({ 
            headless: true,
            executablePath: process.env.CHROMIUM_PATH 
        });
        
        const page = await browser.newPage();
        
        // Ejemplo básico - aquí implementarás tu lógica específica
        await page.goto('https://example.com');
        
        await browser.close();
        
        res.json({ 
            success: true, 
            message: 'Manifiesto generado',
            file_url: 'http://example.com/manifest.pdf'
        });
        
    } catch (error) {
        console.error('Error generando manifiesto:', error);
        res.status(500).json({ 
            success: false, 
            error: error.message 
        });
    }
});

app.listen(PORT, () => {
    console.log(`🎭 Playwright bot running on port ${PORT}`);
});
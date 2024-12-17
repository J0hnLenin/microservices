const express = require('express');
const multer = require('multer');
const path = require('path');
const { Worker } = require('worker_threads')
const now = new Date();
const app = express();
const PORT = 8082;

const storage = multer.memoryStorage();
const upload = multer({ storage: storage });

app.post('/compress', upload.single('image'), async (req, res) => {
    if (!req.file) {
        console.log(getFormattedDate(), 'ERROR: image not found');
        return res.status(400).send('image not found');
    }
    
    try {
        const resultBuffer = await runService(req.file.buffer);
        res.type('image/png');
        
        console.log(getFormattedDate(), 'INFO: image compressed');
        return res.status(200).send(resultBuffer);
        
    } catch (err) {
        console.log(getFormattedDate(), 'ERROR: worker exception', err);
        return res.status(500).send('image compressing exeption');
    }

});

app.listen(PORT, () => {
    console.log(getFormattedDate(), 'start server on port', PORT);
});

function getFormattedDate() {
    const currentDate = new Date();
    return `${currentDate.getFullYear()}-${(currentDate.getMonth() + 1).toString().padStart(2, '0')}-${currentDate.getDate().toString().padStart(2, '0')} ${currentDate.getHours().toString().padStart(2, '0')}:${currentDate.getMinutes().toString().padStart(2, '0')}:${currentDate.getSeconds().toString().padStart(2, '0')}`;
}

function runService(imageBuffer) {
    return new Promise((resolve, reject) => {
        const worker = new Worker('./worker.ts', { workerData: imageBuffer });

        worker.on('message', (result) => {
            const buffer = Buffer.from(result);
            resolve(buffer);
        });
        worker.on('error', reject);
        worker.on('exit', (code) => {
            if (code !== 0) {
                reject(new Error(`Worker stopped with exit code ${code}`));
            }
        });
    });
}


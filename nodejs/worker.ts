const { parentPort, workerData } = require('node:worker_threads');
const sharp = require('sharp');

const processImage = async (imageBuffer) => {
    return await sharp(imageBuffer)
        .resize(28, 28)
        .toBuffer('PNG');
};

processImage(workerData)
    .then(result => parentPort.postMessage(result))
    .catch(err => parentPort.postMessage(err));
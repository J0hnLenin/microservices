import numpy as np
from PIL import Image
from tensorflow import keras
from aiohttp import web
import aiofiles
from datetime import datetime

async def upload_image(request):
    data = await request.post()
    image_file = data['image']
    if image_file:
        file_path = f"./{image_file.filename}"

        async with aiofiles.open(file_path, 'wb') as f:
            await f.write(image_file.file.read())

        try:
            image = Image.open(file_path)
            image_array = np.array(image)
        except:
            timenow = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
            print(f'{timenow} ERROR: invalid image uploaded')
            return web.Response(text="Invalid image uploaded", status=400)

        if image_array.shape[-1] == 3:
            image_array = np.dot(image_array[..., :3], [0.2989, 0.5870, 0.1140])
        elif image_array.shape[-1] == 4:
            image_array = np.dot(image_array[..., :4], [0.2989, 0.5870, 0.1140, 0])

        predicted_digit = predict_digit(image_array)
        timenow = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        print(f'{timenow} INFO: predicted digit: {predicted_digit}')
        return web.Response(text=f'Predicted digit: {predicted_digit}')
    else:
        print(f'{timenow} ERROR: no image uploaded')
        return web.Response(text="No image uploaded", status=400)

def predict_digit(image):
    image = image.reshape(1, 28, 28)
    prediction = model.predict(image)
    return np.argmax(prediction)

async def init_app():
    app = web.Application()
    app.router.add_post('/predict', upload_image)
    return app

if __name__ == '__main__':
    model = keras.models.load_model('digit_model.h5')
    web.run_app(init_app(), port=8083)
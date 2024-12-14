import numpy as np
from PIL import Image
from tensorflow import keras

model = keras.models.load_model('digit_model.h5')

image_path = 'go_image.png'
image = Image.open(image_path).resize((28, 28)) 

image_array = np.array(image)

if image_array.shape[-1] == 3:
    image_array = np.dot(image_array[..., :3], [0.2989, 0.5870, 0.1140])


def predict_digit(image):
    image = image.reshape(1, 28, 28)
    prediction = model.predict(image)
    return np.argmax(prediction)

predicted_digit = predict_digit(image_array)

print(f'Predicted Digit: {predicted_digit}')


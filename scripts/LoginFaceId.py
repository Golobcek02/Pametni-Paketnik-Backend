import glob
import os
import sys

import cv2
import numpy as np
import tensorflow as tf

from Hog import Hog
from Lbp import Lbp
from ViolaJones import process_images

user_id = sys.argv[1].rstrip("\r\n")
folder_path = f"images/{user_id}"
files = os.listdir(folder_path)
image_count = len(files)

vector = []
images = []
for file in glob.glob(f"images/{str(user_id)}/*.*"):
    img = cv2.imread(file)
    img = cv2.rotate(img, cv2.ROTATE_90_COUNTERCLOCKWISE)
    vector.append(img)

vector = process_images(vector)
for img in vector:
    img = cv2.resize(img, (100, 100))
    gray_image = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    lbp_image = Lbp(gray_image)
    hog_descriptor = Hog(gray_image, 8, 2, 9)
    feature_vector = np.concatenate((lbp_image.flatten(), hog_descriptor))
    images.append(feature_vector)

images = np.array(images)
loaded_model = tf.keras.models.load_model("models/" + user_id + ".h5")
labels = np.ones(len(images))

# Predict the labels for the three images
three_predictions = loaded_model.predict(images)
predicted_labels = np.argmax(three_predictions, axis=1)
accuracy = np.mean(predicted_labels == labels)
if accuracy > 0.6:
    print("T")
else:
    print("F")

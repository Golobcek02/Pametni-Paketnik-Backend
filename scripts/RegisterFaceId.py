import glob
import os
import sys

import cv2
import h5py
import numpy as np
import tensorflow as tf
from sklearn.model_selection import train_test_split

from Hog import Hog
from Lbp import Lbp
from ViolaJones import process_images

user_id = sys.argv[1].rstrip("\r\n")

folder_path = f"images/{user_id}"

files = os.listdir(folder_path)

image_count = len(files)

data = []

with h5py.File('models/baseModel.h5', 'r') as f:
    data = f['basemodel'][:]
data = np.array(data)

base_model_images = len(data)

labels = []
for i in range(base_model_images + image_count):
    if i < base_model_images:
        labels.append(0)
    else:
        labels.append(1)

labels = np.array(labels)

VJimg = []
g = 0

for file in glob.glob(f"images/{user_id}/*.*"):
    temp = cv2.imread(file)
    # TOTO LINIJO TESITRAJ KO BO NA APLIKACIJI HOSTANO
    temp = cv2.rotate(temp, cv2.ROTATE_90_COUNTERCLOCKWISE)
    VJimg.append(temp)
VJimg = process_images(VJimg)

if len(VJimg) != image_count:
    print("error")
    iterations = image_count - len(VJimg)
    while image_count != len(VJimg):
        VJimg.append(VJimg[0])

images = []
for img in VJimg:
    if g < image_count:
        img = cv2.resize(img, (100, 100))
        gray_image = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
        lbp_image = Lbp(gray_image)
        hog_descriptor = Hog(gray_image, 8, 2, 9)
        feature_vector = np.concatenate((lbp_image.flatten(), hog_descriptor))
        images.append(feature_vector)
    g = g + 1

images = np.array(images)

data = np.vstack((data, images))

# data = np.resize(data, len(labels))
# data = np.array(data)

# Split data into training and testing sets
X_train, X_test, y_train, y_test = train_test_split(data, labels, train_size=0.9, random_state=80, stratify=labels)
# One-hot encode labels
num_classes = len(np.unique(labels))
y_train_encoded = tf.keras.utils.to_categorical(y_train, num_classes)
y_test_encoded = tf.keras.utils.to_categorical(y_test, num_classes)

# Define neural network architecture
model = tf.keras.models.Sequential([
    tf.keras.layers.Dense(128, activation='relu', input_shape=(X_train.shape[1],)),
    tf.keras.layers.Dropout(0.2),
    tf.keras.layers.Dense(64, activation='relu'),
    tf.keras.layers.Dropout(0.2),
    tf.keras.layers.Dense(num_classes, activation='softmax')  # Update the number of output units
])

# # Compile model
model.compile(optimizer='adam', loss='categorical_crossentropy', metrics=['accuracy'])
#
# # Train model
model.fit(X_train, y_train_encoded, epochs=15, batch_size=32, validation_data=(X_test, y_test_encoded))

# # Evaluate model
test_loss, test_acc = model.evaluate(X_test, y_test_encoded)
model.save("models/" + user_id + ".h5")
print("T")

import glob
import os
import sys

import cv2
import h5py
import numpy as np
import tensorflow as tf
from sklearn.model_selection import train_test_split

from Hog import hog
from Lbp import lbp
from ViolaJones import process_images

# user_id = sys.argv[1].rstrip("\r\n")

# folder_path = f"../images/{user_id}"
#
# files = os.listdir(folder_path)
#
# image_count = len(files)

labels = []
for i in range(530):
    if i < 530:
        labels.append(0)
    else:
        labels.append(1)

labels = np.array(labels)
# with h5py.File('models/basemodel.h5', 'r') as f:
#     data = f['basemodel'][:]
# neke = cv2.imread("../images/nejke.jpg")
# data = np.array(data)
# cv2.imshow("banana", neke)
VJimg = []
g = 0
for file in glob.glob("../images/baseModel/*.*"):
    VJimg.append(cv2.imread(file))
VJimg = np.array(VJimg)
VJimg = process_images(VJimg)

images = []
for img in VJimg:
    # if g < image_count:
    img = cv2.resize(img, (100, 100))
    gray_image = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    lbp_image = lbp(gray_image)
    hog_descriptor = hog(gray_image, 8, 2, 9)
    feature_vector = np.concatenate((lbp_image.flatten(), hog_descriptor))
    images.append(feature_vector)
g = g + 1
images = np.array(images)

with h5py.File("../models/baseModel.h5", 'w') as F:
    dset = F.create_dataset('baseModel', data=images)

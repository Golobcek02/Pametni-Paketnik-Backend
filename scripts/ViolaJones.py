import os

import numpy as np
import cv2


def display_images(image_array):
    for i, image in enumerate(image_array):
        cv2.imshow(f'Image {i + 1}', image)

    cv2.waitKey(0)
    cv2.destroyAllWindows()



def process_images(image_array):
    all_face_images=[]
    for image in image_array:
        grayscale_image = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

        face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades +'haarcascade_frontalface_default.xml')
        detected_faces = face_cascade.detectMultiScale(grayscale_image)


        x, y, w, h = detected_faces[0]
        face_region = image[y:y + h, x:x + w]

        scale_width = w / image.shape[1]
        scale_height = h / image.shape[0]

        new_width = int(image.shape[1] * scale_width)
        new_height = int(image.shape[0] * scale_height)

        resized_image = cv2.resize(face_region, (new_width, new_height), interpolation=cv2.INTER_AREA)


        all_face_images.append(resized_image)

    return np.array(all_face_images)

# Example usage:
# image files =
temp=cv2.imread("nekee.jpg")

processed_images = process_images([temp])
display_images(processed_images)

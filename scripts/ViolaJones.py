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
        grayscale_image = cv.cvtColor(image, cv.COLOR_BGR2GRAY)

        face_cascade = cv.CascadeClassifier('path/to/haarcascade_frontalface_alt.xml')
        detected_faces = face_cascade.detectMultiScale(grayscale_image)

        x, y, w, h = detected_faces[0]  # Assuming there's only one face detected
        face_region = image[y:y + h, x:x + w]

        scale_width = w / image.shape[1]
        scale_height = h / image.shape[0]

        new_width = int(original_image.shape[1] * scale_width)
        new_height = int(original_image.shape[0] * scale_height)

        resized_image = cv.resize(image, (new_width, new_height), interpolation=cv.INTER_AREA)


        all_face_images.append(face_region)

    return np.array(all_face_images)

# Example usage:
# image files =
# processed_images = process_images(image_files)
# display_images(processed_images)

import cv2
import numpy as np


def display_images(image_array):
    for i, image in enumerate(image_array):
        cv2.imshow(f'Image {i + 1}', image)

    cv2.waitKey(0)
    cv2.destroyAllWindows()


def process_images(image_array):
    all_face_images = []
    for image in image_array:
        grayscale_image = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

        face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + 'haarcascade_frontalface_default.xml')
        detected_faces = face_cascade.detectMultiScale(grayscale_image)
        if len(detected_faces) != 0:
            x, y, w, h = detected_faces[0]

            face_region = image[y:y + h, x:x + w]

            scale_width = w / image.shape[1]
            scale_height = h / image.shape[0]

            new_width = int(image.shape[1] * scale_width)
            new_height = int(image.shape[0] * scale_height)

            resized_image = cv2.resize(face_region, (new_width, new_height), interpolation=cv2.INTER_AREA)

            all_face_images.append(resized_image)
        else:
            print("not detected")

    return all_face_images

# Example usage:
# image files =
# temp = cv2.imread("img2.jpg")
# temp = cv2.rotate(temp, cv2.ROTATE_90_COUNTERCLOCKWISE)
# cv2.imshow("bnaa", temp)
# cv2.resize(temp, (1280, 720))
#
# processed_images = process_images([temp])
# display_images(processed_images)

import numpy as np
import cv2


def display_images(image_array):
    for i, image in enumerate(image_array):
        cv2.imshow(f'Image {i + 1}', image)

    cv2.waitKey(0)
    cv2.destroyAllWindows()


def process_images(image_array):
    face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + 'haarcascade_frontalface_default.xml')

    margin = 0.2

    def process_face(face, img):
        (x, y, w, h) = face
        margin_x = int(w * margin)
        margin_y = int(h * margin)

        x_start = max(0, x - margin_x)
        y_start = max(0, y - margin_y)
        x_end = min(img.shape[1], x + w + margin_x)
        y_end = min(img.shape[0], y + h + margin_y)

        face_img = img[y_start:y_end, x_start:x_end].copy()
        face_img = cv2.resize(face_img, (100, 100), interpolation=cv2.INTER_AREA)

        mask = np.zeros(face_img.shape[:2], np.uint8)

        bgdModel = np.zeros((1, 65), np.float64)
        fgdModel = np.zeros((1, 65), np.float64)

        smaller_margin_x = margin_x // 4
        smaller_margin_y = margin_y // 4
        rect = (smaller_margin_x, smaller_margin_y, face_img.shape[1] - 2 * smaller_margin_x,
                face_img.shape[0] - 2 * smaller_margin_y)

        cv2.grabCut(face_img, mask, rect, bgdModel, fgdModel, 5, cv2.GC_INIT_WITH_RECT)

        mask2 = np.where((mask == 2) | (mask == 0), 0, 1).astype('uint8')

        face_img = face_img * mask2[:, :, np.newaxis]

        return face_img

    all_face_images = []

    for image_file in image_array:
        image = cv2.imread(image_file)
        if image is None:
            print(f"Unable to read {image_file}. Skipping...")
            continue

        gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
        faces = face_cascade.detectMultiScale(gray, 1.3, 5)

        face_images = [process_face(face, image) for face in faces]

        if not face_images:
            print(f"No faces found in {image_file}.")
        else:
            all_face_images.extend(face_images)

    return np.array(all_face_images)


# Example usage:
processed_images = process_images(image_files)
display_images(processed_images)
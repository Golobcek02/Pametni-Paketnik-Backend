import numpy as np
import cv2

face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + 'haarcascade_frontalface_default.xml')

img = cv2.imread('neke.jpg')

gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

faces = face_cascade.detectMultiScale(gray, 1.3, 5)

margin = 0.2

face_images = []

for (x, y, w, h) in faces:
    margin_x = int(w * margin)
    margin_y = int(h * margin)

    x_start = max(0, x - margin_x)
    y_start = max(0, y - margin_y)
    x_end = min(img.shape[1], x + w + margin_x)
    y_end = min(img.shape[0], y + h + margin_y)

    face_img = img[y_start:y_end, x_start:x_end].copy()

    mask = np.zeros(face_img.shape[:2], np.uint8)

    bgdModel = np.zeros((1, 65), np.float64)
    fgdModel = np.zeros((1, 65), np.float64)

    smaller_margin_x = margin_x // 4
    smaller_margin_y = margin_y // 4
    rect = (smaller_margin_x, smaller_margin_y, face_img.shape[1] - 2*smaller_margin_x, face_img.shape[0] - 2*smaller_margin_y)

    cv2.grabCut(face_img, mask, rect, bgdModel, fgdModel, 5, cv2.GC_INIT_WITH_RECT)

    mask2 = np.where((mask == 2) | (mask == 0), 0, 1).astype('uint8')

    face_img = face_img * mask2[:, :, np.newaxis]

    face_images.append(face_img)

for i, face_img in enumerate(face_images):
    cv2.imshow(f'Face {i+1}', face_img)

cv2.waitKey()
cv2.destroyAllWindows()

import numpy as np

def Lbp(image):

    height, width = image.shape
    lbp_image = np.zeros((height, width), dtype=np.uint8)
    image = np.pad(image, pad_width=1)

    for ih in range(1, height - 1):
        for iw in range(1, width - 1):
            center = image[ih, iw]
            neighbors = [
                image[ih - 1, iw - 1],
                image[ih - 1, iw],
                image[ih - 1, iw + 1],
                image[ih, iw - 1],
                image[ih, iw + 1],
                image[ih + 1, iw - 1],
                image[ih + 1, iw],
                image[ih + 1, iw + 1],
            ]

            binary_pattern = (np.array(neighbors) >= center).astype(int)
            powers_of_two = np.array([1, 2, 4, 8, 16, 32, 64, 128], dtype=np.uint8)
            lbp_value = np.sum(binary_pattern * powers_of_two)

            lbp_image[ih, iw] = lbp_value

    return lbp_image
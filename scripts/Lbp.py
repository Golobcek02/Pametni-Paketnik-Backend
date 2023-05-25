import numpy as np

def lbp(image):

    height, width = image.shape
    lbp_image = np.zeros((height, width), dtype=np.uint8)

    for ih in range(1, height - 1):
        for iw in range(1, width - 1):
           

    return lbp_image
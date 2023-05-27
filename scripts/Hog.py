import numpy as np

def Hog(gray_image, cell_size, block_size, num_bins):
    height, width = gray_image.shape
    gx = np.zeros_like(gray_image, dtype=np.float32)
    gy = np.zeros_like(gray_image, dtype=np.float32)

    gx[:, :-1] = np.diff(gray_image, n=1, axis=1)
    gy[:-1, :] = np.diff(gray_image, n=1, axis=0)

    gradient_magnitude = np.sqrt(gx ** 2 + gy ** 2)
    gradient_orientation = (np.arctan2(gy, gx) * (180.0 / np.pi) + 180.0) % 180.0

    num_cells_x = width // cell_size
    num_cells_y = height // cell_size

    histogram = np.zeros((num_cells_y, num_cells_x, num_bins))

    for y in range(num_cells_y):
        for x in range(num_cells_x):
            cell_magnitude = gradient_magnitude[y * cell_size:(y + 1) * cell_size,
                             x * cell_size:(x + 1) * cell_size]
            cell_orientation = gradient_orientation[y * cell_size:(y + 1) * cell_size,
                               x * cell_size:(x + 1) * cell_size]
            
            hist = np.histogram(cell_orientation, bins=num_bins, range=(0, 180),
                                weights=cell_magnitude)[0]

            histogram[y, x, :] = hist / np.sqrt(np.sum(hist ** 2) + 1e-6)

    hog_descriptor = np.zeros((num_cells_y - block_size + 1, num_cells_x - block_size + 1,
                               block_size * block_size * num_bins))

    for y in range(num_cells_y - block_size + 1):
        for x in range(num_cells_x - block_size + 1):
            block_histogram = histogram[y:y + block_size, x:x + block_size, :].flatten()
            hog_descriptor[y, x, :] = block_histogram / np.sqrt(np.sum(block_histogram ** 2) + 1e-6)

    return hog_descriptor.flatten()
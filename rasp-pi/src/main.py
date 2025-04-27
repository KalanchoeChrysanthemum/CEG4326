import cv2
import numpy as np

# Load the cropped 4x4 LED image
image = cv2.imread('../images/input/a2.jpg', cv2.IMREAD_GRAYSCALE)
cropped = image[500:1350, 600:1600]
#cropped = cv2.GaussianBlur(cropped, (11, 11), 0)
_, cropped = cv2.threshold(cropped, 220, 255, cv2.THRESH_BINARY)

# Grid size
rows, cols = 4, 4

# Image dimensions
height, width = cropped.shape[:2]

# Size of each cell
cell_h = height // rows
cell_w = width // cols

# Draw rectangles over each cell

binary = ""

for i in range(rows):
    for j in range(cols):
        x = j * cell_w
        y = i * cell_h
        roi = cropped[y:y+cell_h, x:x+cell_w]

        avg_brightness = int(np.mean(roi))
        status = '1' if avg_brightness > 20 else '0'
        binary += status
        
        color = (0, 255, 0) if avg_brightness > 150 else (0, 0, 255)
        cv2.rectangle(cropped, (x, y), (x + cell_w, y + cell_h), (255, 255, 255), 2)
        text = str(avg_brightness)
        text_size = cv2.getTextSize(text, cv2.FONT_HERSHEY_SIMPLEX, 0.5, 1)[0]
        text_x = x + (cell_w - text_size[0]) // 2
        text_y = y + (cell_h + text_size[1]) // 2
        cv2.putText(cropped, text, (text_x, text_y), cv2.FONT_HERSHEY_SIMPLEX, 0.5, color, 1)
    binary += '\n'

# Save the output
cv2.imwrite('../images/output/processed_a3.jpg', cropped)
print(binary)


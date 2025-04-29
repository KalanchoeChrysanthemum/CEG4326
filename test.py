import cv2
import numpy as np

# Load the image
image_path = 'uh.jpg'
img = cv2.imread(image_path)
gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

threshold = 200
_, binary = cv2.threshold(gray, threshold, 255, cv2.THRESH_BINARY)

# Get image dimensions
height, width = gray.shape

num_bulbs = 6

# Calculate the width of each ROI
roi_width = width // num_bulbs

# Define the height of the ROI and shift it downward to cover the lights
# The lights are roughly in the middle of the image vertically, so we'll start the ROI lower
roi_height = int(height * 0.4)  # Keep the height as 40% of the image height
y_start = int(height * 0.3)     # Start the ROI 30% down from the top of the image (adjust this as needed)
y_end = y_start + roi_height    # End of the ROI vertically

# Ensure y_end doesn't exceed image height
y_end = min(y_end, height)

# Create a copy of the image for visualization
vis = binary.copy()
binary_output = ""

# Process each bulb
for i in range(num_bulbs):
    # Calculate the x-coordinates for the ROI
    x_start = i * roi_width
    x_end = x_start + roi_width
    
    # Ensure x_end doesn't exceed image width
    x_end = min(x_end, width)
    
    # Extract the ROI (shifted downward)
    roi = gray[y_start:y_end, x_start:x_end]
    
    # Calculate average brightness
    avg_brightness = np.mean(roi)
    status = '1' if avg_brightness > 150 else '0'
    binary_output += status

    # Draw the ROI rectangle and label
    color = (0, 255, 0) if status == '1' else (0, 0, 255)
    cv2.rectangle(vis, (x_start, y_start), (x_end, y_end), color, 3)
    cv2.putText(vis, f'{status} ({int(avg_brightness)})', (x_start, y_end + 30), 
                cv2.FONT_HERSHEY_SIMPLEX, 0.7, color, 2)

# Save the output image
cv2.imwrite('uhh.jpg', vis)
print(f"Binary output: {binary_output}")

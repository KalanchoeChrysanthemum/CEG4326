import cv2

image = cv2.imread('nwero.jpg')

gray_scale = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

cv2.imwrite('gray_nwero.jpg', gray_scale)

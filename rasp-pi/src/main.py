import cv2

INPUT_IMG_PATH = '../images/input/nwero.jpg'
OUTPUT_IMG_PATH = '../images/output/processed_nwero.jpg'

def process_img(base_img):
    return cv2.adaptiveThreshold(base_img, 255, cv2.ADAPTIVE_THRESH_GAUSSIAN_C, cv2.THRESH_BINARY, 11, 2)

def main():
    img = cv2.imread(INPUT_IMG_PATH, cv2.IMREAD_GRAYSCALE)
    cv2.imwrite("../images/output/gray_nwero.jpg", img)
    cv2.imwrite(OUTPUT_IMG_PATH, process_img(img))

if __name__ == '__main__':
    main()

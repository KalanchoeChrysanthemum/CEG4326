import cv2

INPUT_IMG_PATH = 'nwero.jpg'
OUTPUT_IMG_PATH = 'gray_nwero.jpg'

def process_img(base_img):
    gray_img = cv2.cvtColor(base_img, cv2.COLOR_BGR2GRAY) # Convert to gray scale
    return gray_img

def main():
    base_img = cv2.imread(INPUT_IMG_PATH)
    cv2.imwrite(OUTPUT_IMG_PATH, process_img(base_img))

if __name__ == '__main__':
    main()

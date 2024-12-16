from colorama import init
from termcolor import colored
import cv2
from point import Point

def img(name):
    """
    Returns an OpenCV image
    """

    return cv2.imread(name, 0)

def square_img(image):
    """
    Returns all the tools to consider the image <image> as a square
    """
    print("Squaring the image...", end=" ")
    height, width = image.shape[0], image.shape[1]
    length = min(height, width)
    center = Point(image.shape[0] / 2, image.shape[1] / 2)
    radius = length / 2 - 1 / 2
    print(colored("DONE", "green"))
    return center, radius, width, height, length

def show_img(image):
    """
    Shows the image <image>
    """

    cv2.imshow("image", image)
    cv2.waitKey(0)
    cv2.destroyAllWindows()

def circle_img(width, height, radius, length, center, image):
    """
    Takes a circle from the middle of image <image>
    """

    print("Extracting a circle from given image...", end=" ")

    for y in range(height):
        for x in range(width):
            if Point(x,y).distance(center) > radius:
                image[y,x] = 0xFF

    print(colored("DONE", "green"))

def drawline(image, origin, end):
    cv2.line(image, (origin.x * 25, origin.y * 25), (end.x * 25, end.y * 25), color=0, thickness=4, lineType=8)

def shrink_image(image, target_height, target_width):
    cv2.resize(image, (target_width, target_height))

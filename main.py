#!/usr/bin/env python3

from point import Point
from image import *
from point_generation import *

from cv2 import cv2
from termcolor import colored
import numpy as np
import os
import math
import time

IMG = "./images/zuko3.png"
DECOMPOSITION = False
NUMBER_LINES = 10000
NUMBER_POINTS = 360

def get_data(dictionnary, pt, pt2):
    """
    Since all_lines doesn't have a symmetric data for <pt> <pt2>, picks the good one
    """

    a, b = pt, pt2
    if pt.id > pt2.id:
        a, b = b, a

    if a not in dictionnary:
        distance = int(a.distance(b))
        xs = np.linspace(a.x, b.x, distance, dtype=int)
        ys = np.linspace(a.y, b.y, distance, dtype=int)
        weight = 1
        dictionnary[a] = {b: (xs, ys, distance, weight)}

    elif b not in dictionnary[a]:
        distance = int(a.distance(b))
        xs = np.linspace(a.x, b.x, distance, dtype=int)
        ys = np.linspace(a.y, b.y, distance, dtype=int)
        weight = 1
        dictionnary[a][b] = (xs, ys, distance, weight)

    return dictionnary[a][b]

def main():
    tic = time.time()
    image_name = IMG
    number_points = NUMBER_POINTS
    number_lines = NUMBER_LINES

    image = img(image_name)
    center, radius, width, height, length = square_img(image)
    circle_img(width, height, radius, length, center, image)


    coords = generate_points(center, radius, number_points)
    lines = {}
    error = np.ones(image.shape) * 0xFF - image.copy()
    result = np.ones((image.shape[0] * 25, image.shape[1] * 25), np.uint8) * 0xFF
    mask = np.zeros(image.shape, np.float64)

    order_points = [coords[0]]
    last = coords[0]

    for l in range(number_lines):
        evolution = "Drawing {} lines... {}".format(number_lines, colored("({}/{})".format(l, number_lines), "yellow"))
        print(evolution, end="\r")

        ressemblance = -math.inf
        best_choice = last

        # Find the line which will lower the error the most

        for candidat in coords:
            xs, ys, _, weight = get_data(lines, last, candidat)

            line_ressemblance = np.sum(error[ys, xs]) * weight
            if line_ressemblance > ressemblance:
                ressemblance = line_ressemblance
                best_choice = candidat

        order_points.append(best_choice)
        xs, ys, dist, weight = get_data(lines, best_choice, last)
        weight *= 15
        if dist == 0:
           break

        mask.fill(0)
        mask[ys, xs] = weight
        error -= mask
        error.clip(0, 255)

        drawline(result, last, best_choice)
        last = best_choice
        if DECOMPOSITION:
            #################################################
            name = str(l) + ".png"
            cv2.imwrite(name, result)
            #################################################

    print(" "*len(evolution), end="\r")
    print("Drawing {} lines... {}".format(number_lines, colored("DONE", "green")))
    print("Creating image... ", end="")
    name = os.path.splitext(IMG)[0] + "-string-art.png"
    cv2.imwrite(name, result)
    print(colored("DONE", "green"))

    print("Program ended successfully in {} seconds".format(time.time() - tic))


if __name__ == "__main__":
    main()
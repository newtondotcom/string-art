import math
from point import Point
from colorama import init
from termcolor import colored

def generate_points(center, circ_radius, n=360):
    """
    Generates a circle with <n> points
    """

    
    if n < 36:
        n = 360

    coords = []
    for pt in range(n):
        evolution = "Generating {} points in circle... {}".format(n, colored("({}/{})".format(pt, n), "yellow"))
        print(evolution, end="\r")
        angle = (2 * math.pi * pt) / n
        x = math.floor(center.x + circ_radius * math.cos(angle))
        y = math.floor(center.y + circ_radius * math.sin(angle))
        coords.append(Point(x, y, pt))
    print(" "*len(evolution), end="\r")
    print("Generating {} points in circle... {}".format(n, colored("DONE", "green")))

    return coords

import math

class Point:
    def __init__(self, x, y, id=-1):
        self.x = x
        self.y = y
        self.id = id

    def distance(self, other):
        """
        Returns the distance between two points
        """

        return math.sqrt((other.x - self.x) ** 2 + (other.y - self.y) ** 2)

    def __str__(self):
        """
        Representation of a point
        """

        return '{} ({},{})'.format(self.id, self.x, self.y)
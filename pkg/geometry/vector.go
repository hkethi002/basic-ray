package geometry

import "math"

func DotProduct(vectorA, vectorB Vector) float64 {
	return (vectorA[0] * vectorB[0]) + (vectorA[1] * vectorB[1]) + (vectorA[2] * vectorB[2])
}

func CrossProduct(vectorA, vectorB Vector) Vector {
	return Vector{
		vectorA[1]*vectorB[2] - vectorA[2]*vectorB[1],
		vectorA[2]*vectorB[0] - vectorA[0]*vectorB[2],
		vectorA[0]*vectorB[1] - vectorA[1]*vectorB[0],
	}
}

func ScalarProduct(vector Vector, scalar float64) Vector {
	return Vector{vector[0] * scalar, vector[1] * scalar, vector[2] * scalar}
}

func Magnitude(vector Vector) float64 {
	return math.Sqrt(math.Pow(vector[0], 2) + math.Pow(vector[1], 2) + math.Pow(vector[2], 2))
}

func Normalize(vector Vector) Vector {
	magnitude := Magnitude(vector)
	return ScalarProduct(vector, 1/magnitude)
}

func Subtract(vectorA, vectorB Vector) Vector {
	return Vector{vectorA[0] - vectorB[0], vectorA[1] - vectorB[1], vectorA[2] - vectorB[2]}
}

func Add(vectorA, vectorB Vector) Vector {
	return Vector{vectorA[0] + vectorB[0], vectorA[1] + vectorB[1], vectorA[2] + vectorB[2]}
}

// origin at B to A
func CreateVector(pointA, pointB Point) Vector {
	return Vector{pointA[0] - pointB[0], pointA[1] - pointB[1], pointA[2] - pointB[2]}
}

func Distance(pointA, pointB Point) float64 {
	return Magnitude(CreateVector(pointA, pointB))
}

func Translate(point Point, vector Vector) Point {
	return Point{point[0] + vector[0], point[1] + vector[1], point[2] + vector[2]}
}

func (vector *Vector) IsNormal() bool {
	return Magnitude(*vector) == 1
}

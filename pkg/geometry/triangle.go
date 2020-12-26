package geometry

func (triangle *Triangle) GetNormal() Vector {
	edge1 := CreateVector(triangle.vertex1, triangle.vertex0)
	edge2 := CreateVector(triangle.vertex2, triangle.vertex0)

	normal := CrossProduct(edge1, edge2)
	return normal
}

func getLinePoint(ray *Ray, d float64) Point {
	addPoint := ScalarProduct(ray.Vector, d)
	return Point{
		ray.Origin[0] + addPoint[0],
		ray.Origin[1] + addPoint[1],
		ray.Origin[2] + addPoint[2],
	}
}

func GetIntersection(ray *Ray, triangle *Triangle) *Point {
	epsilon := 0.0000001
	normal := triangle.GetNormal()
	cos := DotProduct(ray.Vector, normal)
	if cos == 0 {
		return nil
	}

	d := DotProduct(CreateVector(triangle.vertex0, ray.Origin), normal) / cos
	potentialPoint := getLinePoint(ray, d)

	var subTriangle *Triangle
	subTriangleArea := 0.0

	subTriangle = &Triangle{vertex0: triangle.vertex0, vertex1: triangle.vertex1, vertex2: potentialPoint}
	subTriangleArea += subTriangle.GetArea()

	subTriangle.vertex0 = triangle.vertex2
	subTriangleArea += subTriangle.GetArea()

	subTriangle.vertex1 = triangle.vertex0
	subTriangleArea += subTriangle.GetArea()

	if subTriangleArea-triangle.GetArea() < epsilon {
		return &potentialPoint
	}
	return nil
}

func (triangle *Triangle) GetArea() float64 {
	return Magnitude(triangle.GetNormal()) / 2.0
}

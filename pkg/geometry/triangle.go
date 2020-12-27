package geometry

func (triangle *Triangle) GetNormal() Vector {
	if triangle.Normal != nil {
		return *triangle.Normal
	}
	edge1 := CreateVector(triangle.Vertex1, triangle.Vertex0)
	edge2 := CreateVector(triangle.Vertex2, triangle.Vertex0)

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

	d := DotProduct(CreateVector(triangle.Vertex0, ray.Origin), normal) / cos
	if d < epsilon {
		return nil // I think this means point is in wrong direction
	}

	potentialPoint := getLinePoint(ray, d)

	var subTriangle *Triangle
	subTriangleArea := 0.0

	subTriangle = &Triangle{Vertex0: triangle.Vertex0, Vertex1: triangle.Vertex1, Vertex2: potentialPoint}
	subTriangleArea += subTriangle.GetArea()

	subTriangle.Vertex0 = triangle.Vertex2
	subTriangleArea += subTriangle.GetArea()

	subTriangle.Vertex1 = triangle.Vertex0
	subTriangleArea += subTriangle.GetArea()

	if subTriangleArea-triangle.GetArea() < epsilon {
		return &potentialPoint
	}
	return nil
}

func (triangle *Triangle) GetArea() float64 {
	return Magnitude(triangle.GetNormal()) / 2.0
}

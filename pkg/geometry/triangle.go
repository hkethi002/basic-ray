package geometry

func (triangle *Triangle) GetNormal() Vector {
	if Magnitude(triangle.Normal) > 0 {
		return triangle.Normal
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
	edge1 := CreateVector(triangle.Vertex1, triangle.Vertex0)
	edge2 := CreateVector(triangle.Vertex2, triangle.Vertex0)
	return Magnitude(CrossProduct(edge1, edge2)) / 2.0
}

func TriangulatePolygon(points []Point, normal Vector, vertexNormals []Vector, texture *TextureProperties) []*Triangle {
	triangles := make([]*Triangle, len(points)-2)
	for i := 0; i < len(points)-2; i++ {
		triangle := &Triangle{
			Vertex0:            points[0],
			Vertex1:            points[i+1],
			Vertex2:            points[i+2],
			Normal:             normal,
			VertexNormals:      vertexNormals,
			DiffuseAlbedo:      texture.DiffuseAlbedo,
			SpecularAlbedo:     texture.SpecularAlbedo,
			TranslucenseAlbedo: texture.TranslucenseAlbedo,
			MaterialType:       texture.MaterialType,
		}
		triangles[i] = triangle
	}
	return triangles
}

func GetFaceVertexes(face []int, vertexes []Point) []Point {
	points := make([]Point, len(face))
	for faceIndex, vertexIndex := range face {
		points[faceIndex] = vertexes[vertexIndex]
	}
	return points
}

func TriangulateObject(object *Object) []*Triangle {
	var triangles []*Triangle

	for i, face := range object.Faces {
		var vertexNormals []Vector
		points := GetFaceVertexes(face, object.Vertexes)
		texture := object.Textures[object.TextureMap[i]]
		if len(object.VertexNormals) > 0 {
			vertexNormals = make([]Vector, len(points))
			for p := 0; p < len(points); p++ {
				vertexNormals[p] = Normalize(object.VertexNormals[face[p]])
			}
		}
		normal := Normalize(object.Normals[i])

		triangles = append(
			triangles,
			TriangulatePolygon(points, normal, vertexNormals, &texture)...,
		)
	}
	return triangles
}

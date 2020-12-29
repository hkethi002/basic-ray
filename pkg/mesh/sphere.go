package mesh

/*
Package for making meshes, outputs should be json Object
structs from geometry which can be exported as json
*/
import (
	geometry "basic-ray/pkg/geometry"
)

const ICO_EDGE_LENGTH = 0.9510565163 // sin(2*pi/5)

type Sphere struct {
	Radius float64
	Origin geometry.Point
}

func (sphere *Sphere) CreateMesh(refinement int) *Mesh {
	mesh := sphere.createBaseIcosahedron()
	for i := 0; i < refinement; i++ {
		mesh = RefineMesh(sphere, mesh)
	}
	return mesh

}

func GetMidPoint(sphere *Sphere, mesh *Mesh, faceIndexA, faceIndexB int, cache *map[[2]int]int) int {
	// Check cache
	vertexIndex, present := (*cache)[[2]int{faceIndexA, faceIndexB}]
	if present == true {
		return vertexIndex
	}
	vertexIndex, present = (*cache)[[2]int{faceIndexB, faceIndexA}]
	if present == true {
		return vertexIndex
	}

	// Calculate midPoint vertex if not in cache
	vertexA := mesh.Vertexes[faceIndexA]
	vertexB := mesh.Vertexes[faceIndexB]

	translate := geometry.ScalarProduct(geometry.CreateVector(vertexB, vertexA), .5)
	linePoint := geometry.Translate(vertexA, translate)

	vector := geometry.ScalarProduct(geometry.Normalize(geometry.CreateVector(linePoint, sphere.Origin)), sphere.Radius)
	midVertex := geometry.Translate(sphere.Origin, vector)

	// Add midPoint vertex to the mesh and cache
	mesh.Vertexes = append(mesh.Vertexes, midVertex)
	vertexIndex = len(mesh.Vertexes) - 1
	(*cache)[[2]int{faceIndexA, faceIndexB}] = vertexIndex
	return vertexIndex
}

func RefineMesh(sphere *Sphere, mesh *Mesh) *Mesh {
	refinedMesh := Mesh{Vertexes: mesh.Vertexes}
	midPointCache := make(map[[2]int]int)
	faces := make([][]int, 0)
	for _, face := range mesh.Faces {
		newVertexes := make([]int, 3)
		for edge := 0; edge < 3; edge++ {
			newVertexes[edge] = GetMidPoint(sphere, &refinedMesh, face[edge], face[(edge+1)%3], &midPointCache)
		}

		faces = append(faces, []int{face[0], newVertexes[0], newVertexes[2]})
		faces = append(faces, []int{face[1], newVertexes[1], newVertexes[0]})
		faces = append(faces, []int{face[2], newVertexes[2], newVertexes[1]})
		faces = append(faces, []int{newVertexes[0], newVertexes[1], newVertexes[2]})

	}

	normals := make([]geometry.Vector, len(faces))

	for i, face := range faces {
		edge1 := geometry.CreateVector(refinedMesh.Vertexes[face[0]], refinedMesh.Vertexes[face[1]])
		edge2 := geometry.CreateVector(refinedMesh.Vertexes[face[0]], refinedMesh.Vertexes[face[2]])
		normals[i] = geometry.Normalize(geometry.CrossProduct(edge1, edge2))
	}
	refinedMesh.Faces = faces
	refinedMesh.Normals = normals
	return &refinedMesh
}

func (sphere *Sphere) createBaseIcosahedron() *Mesh {
	sideLength := sphere.Radius * ICO_EDGE_LENGTH / 2.0
	translate := geometry.CreateVector(sphere.Origin, geometry.Point{0, 0, 0})
	vertexes := []geometry.Point{
		geometry.Translate(geometry.Point{2 * sideLength, sideLength, 0}, translate),
		geometry.Translate(geometry.Point{-2 * sideLength, sideLength, 0}, translate),
		geometry.Translate(geometry.Point{2 * sideLength, -1 * sideLength, 0}, translate),
		geometry.Translate(geometry.Point{-2 * sideLength, -1 * sideLength, 0}, translate),

		geometry.Translate(geometry.Point{0, 2 * sideLength, sideLength}, translate),
		geometry.Translate(geometry.Point{0, -2 * sideLength, sideLength}, translate),
		geometry.Translate(geometry.Point{0, 2 * sideLength, -1 * sideLength}, translate),
		geometry.Translate(geometry.Point{0, -2 * sideLength, -1 * sideLength}, translate),

		geometry.Translate(geometry.Point{sideLength, 0, 2 * sideLength}, translate),
		geometry.Translate(geometry.Point{sideLength, 0, -2 * sideLength}, translate),
		geometry.Translate(geometry.Point{-1 * sideLength, 0, 2 * sideLength}, translate),
		geometry.Translate(geometry.Point{-1 * sideLength, 0, -2 * sideLength}, translate),
	}

	faces := [][]int{
		[]int{0, 6, 4},
		[]int{0, 8, 2},
		[]int{0, 2, 9},
		[]int{0, 9, 6},
		[]int{0, 4, 8},

		[]int{6, 1, 4},
		[]int{4, 10, 8},
		[]int{8, 5, 2},
		[]int{9, 2, 7},
		[]int{0, 11, 6},

		[]int{3, 10, 1},
		[]int{3, 1, 11},
		[]int{3, 5, 10},
		[]int{3, 7, 5},
		[]int{3, 11, 7},

		[]int{1, 10, 4},
		[]int{10, 5, 8},
		[]int{5, 7, 2},
		[]int{7, 11, 9},
		[]int{11, 1, 6},
	}

	normals := make([]geometry.Vector, 20)

	for i, face := range faces {
		edge1 := geometry.CreateVector(vertexes[face[0]], vertexes[face[1]])
		edge2 := geometry.CreateVector(vertexes[face[0]], vertexes[face[2]])
		normals[i] = geometry.Normalize(geometry.CrossProduct(edge1, edge2))
	}

	return &Mesh{Vertexes: vertexes, Faces: faces, Normals: normals}
}

package main

import (
	"math"

	"github.com/thegreatdaniad/go-tracer/obj_parser"
	"github.com/ungerik/go3d/mat3"
	"github.com/ungerik/go3d/vec3"
)

type Geometry interface {
	GetGeometryData() GeometryData
	SetMaterial(material Material)
}

type GeometryData struct {
	Vertices           []vec3.T
	Faces              []Face
	Material           Material
	TextureCoordinates []TextureCoordinate
	Normals            []Normal
}

type Normal struct {
	X, Y, Z float32
}

func (n Normal) ToVec3() vec3.T {
	return vec3.T{float32(n.X), float32(n.Y), float32(n.Z)}
}

type TextureCoordinate struct {
	U, V float64
}

type Face struct {
	VertexIndices            []int
	TextureCoordinateIndices []int
	NormalIndices            []int
	Material                 Material
}

type Obj struct {
	Vertices           []vec3.T
	TextureCoordinates []TextureCoordinate
	Normals            []Normal
	Faces              []Face
	Origin             vec3.T
}

func (o *Obj) Rotate(degX, degY, degZ float64) {
	// Convert degrees to radians
	radX := degX * math.Pi / 180
	radY := degY * math.Pi / 180
	radZ := degZ * math.Pi / 180

	// Rotation matrices for X, Y, Z axes
	rotX := mat3.T{
		vec3.T{1, 0, 0},
		vec3.T{0, float32(math.Cos(radX)), -float32(math.Sin(radX))},
		vec3.T{0, float32(math.Sin(radX)), float32(math.Cos(radX))},
	}
	rotY := mat3.T{
		vec3.T{float32(math.Cos(radY)), 0, float32(math.Sin(radY))},
		vec3.T{0, 1, 0},
		vec3.T{-float32(math.Sin(radY)), 0, float32(math.Cos(radY))},
	}
	rotZ := mat3.T{
		vec3.T{float32(math.Cos(radZ)), -float32(math.Sin(radZ)), 0},
		vec3.T{float32(math.Sin(radZ)), float32(math.Cos(radZ)), 0},
		vec3.T{0, 0, 1},
	}

	for i, vertex := range o.Vertices {
		// Translate vertex to origin
		translated := vec3.Sub(&vertex, &o.Origin)

		rotated := rotX.MulVec3(&translated)
		rotated = rotY.MulVec3(&rotated)
		rotated = rotZ.MulVec3(&rotated)

		// Translate vertex back
		o.Vertices[i] = vec3.Add(&rotated, &o.Origin)
	}

	for i, normal := range o.Normals {
		normalVec := vec3.T{float32(normal.X), float32(normal.Y), float32(normal.Z)}

		// Apply rotations
		rotatedNormal := rotX.MulVec3(&normalVec)
		rotatedNormal = rotY.MulVec3(&rotatedNormal)
		rotatedNormal = rotZ.MulVec3(&rotatedNormal)

		o.Normals[i] = Normal{(rotatedNormal[0]), (rotatedNormal[1]), (rotatedNormal[2])}
	}
}

func (o *Obj) GetGeometryData() GeometryData {
	return GeometryData{
		Vertices: o.Vertices,
		Faces:    o.Faces,
		Normals:  o.Normals,
	}
}
func (o *Obj) SetMaterial(material Material) {
	for i := range o.Faces {
		o.Faces[i].Material = material
	}
}

func ParseObjFile(filename string) (*Obj, error) {
	o, err := obj_parser.ParseObjFile(filename)
	if err != nil {
		return nil, err
	}

	newObj := &Obj{
		Vertices:           o.Vertices,
		Normals:            make([]Normal, len(o.Vertices)), // Allocate space for vertex normals
		TextureCoordinates: make([]TextureCoordinate, len(o.TextureCoordinates)),
		Faces:              make([]Face, len(o.Faces)),
	}

	// Copy texture coordinates
	for i, tc := range o.TextureCoordinates {
		newObj.TextureCoordinates[i] = TextureCoordinate{
			U: tc.U,
			V: tc.V,
		}
	}

	// Initialize a slice to hold the accumulated normals for each vertex
	accumulatedNormals := make([]vec3.T, len(o.Vertices))

	// Initialize a counter for each vertex to count how many faces share the vertex
	sharedFaces := make([]int, len(o.Vertices))

	// Iterate over all faces to calculate face normals and accumulate them
	for _, f := range o.Faces {
		V1 := newObj.Vertices[f.VertexIndices[0]]
		V2 := newObj.Vertices[f.VertexIndices[1]]
		V3 := newObj.Vertices[f.VertexIndices[2]]

		// Calculate the face normal
		normal := calculateFaceNormal(V1, V2, V3)

		// Accumulate the face normal to each vertex normal and increment the shared face count
		for _, vIdx := range f.VertexIndices {
			accumulatedNormals[vIdx] = *accumulatedNormals[vIdx].Add(&normal)
			sharedFaces[vIdx]++
		}
	}

	// Average the accumulated normals by the number of shared faces and assign them to newObj.Normals
	for i := range newObj.Normals {
		if sharedFaces[i] > 0 {
			accumulatedNormals[i].Scale(1.0 / float32(sharedFaces[i]))
			accumulatedNormals[i].Normalize() // Normalize the averaged normal
			newObj.Normals[i] = Normal{X: accumulatedNormals[i][0], Y: accumulatedNormals[i][1], Z: accumulatedNormals[i][2]}
		}
	}

	// Assign faces to newObj
	for i, f := range o.Faces {
		newFace := Face{
			VertexIndices:            make([]int, len(f.VertexIndices)),
			TextureCoordinateIndices: make([]int, len(f.TextureCoordinateIndices)),
			NormalIndices:            make([]int, len(f.VertexIndices)), // Initialize with the same length as VertexIndices
		}
		copy(newFace.VertexIndices, f.VertexIndices)
		copy(newFace.TextureCoordinateIndices, f.TextureCoordinateIndices)

		// Now, assign the corresponding vertex normal indices to the face
		for j := range newFace.VertexIndices {
			newFace.NormalIndices[j] = newFace.VertexIndices[j]
		}
		newObj.Faces[i] = newFace
	}

	return newObj, nil
}

// calculateFaceNormal assumes that the vertices are in counter-clockwise order
func calculateFaceNormal(V1, V2, V3 vec3.T) vec3.T {
	edge1 := vec3.Sub(&V2, &V1)
	edge2 := vec3.Sub(&V3, &V1)
	normal := vec3.Cross(&edge1, &edge2)
	normal.Normalize()
	return normal
}

// func calculateFaceNormal(vertices []vec3.T, vertexIndices []int) [3]float32 {
// 	if len(vertexIndices) < 3 {
// 		return [3]float32{0, 0, 0}
// 	}

// 	v0 := vertices[vertexIndices[0]]
// 	v1 := vertices[vertexIndices[1]]
// 	v2 := vertices[vertexIndices[2]]

// 	edge1 := vec3.Sub(&v1, &v0)
// 	edge2 := vec3.Sub(&v2, &v0)

// 	normalVec := vec3.Cross(&edge1, &edge2)

// 	normalVec.Normalize()

//		return [3]float32{normalVec[0], normalVec[1], normalVec[2]}
//	}
func (f Face) Intersects(ray Ray, vertices []vec3.T) (bool, float32, vec3.T) {
	var intersectionPoint vec3.T // Declare the intersection point variable

	if len(f.VertexIndices) < 3 {
		return false, 0, intersectionPoint
	}

	v0 := vertices[f.VertexIndices[0]]
	v1 := vertices[f.VertexIndices[1]]
	v2 := vertices[f.VertexIndices[2]]

	// Edge vectors
	e1 := vec3.Sub(&v1, &v0)
	e2 := vec3.Sub(&v2, &v0)

	// Begin calculating determinant - also used to calculate u parameter
	pvec := vec3.Cross(&ray.Direction, &e2)

	// If determinant is near zero, ray lies in plane of triangle
	det := vec3.Dot(&e1, &pvec)

	// NOT CULLING
	if det > -float32(math.SmallestNonzeroFloat32) && det < float32(math.SmallestNonzeroFloat32) {
		return false, 0, intersectionPoint
	}
	invDet := 1.0 / det

	// Calculate distance from V0 to ray origin
	tvec := vec3.Sub(&ray.Origin, &v0)

	// Calculate u parameter and test bound
	u := vec3.Dot(&tvec, &pvec) * invDet

	// The intersection lies outside of the triangle
	if u < 0.0 || u > 1.0 {
		return false, 0, intersectionPoint
	}

	// Prepare to test v parameter
	qvec := vec3.Cross(&tvec, &e1)

	// Calculate V parameter and test bound
	v := vec3.Dot(&ray.Direction, &qvec) * invDet

	// The intersection lies outside of the triangle
	if v < 0.0 || u+v > 1.0 {
		return false, 0, intersectionPoint
	}

	t := vec3.Dot(&e2, &qvec) * invDet

	// Calculate the exact point of intersection
	if t > float32(math.SmallestNonzeroFloat32) {
		rs := ray.Direction.Scaled(t)
		ip := ray.Origin.Add(&rs)
		intersectionPoint = *ip
		return true, t, intersectionPoint
	}

	return false, 0, intersectionPoint
}

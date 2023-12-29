package main

import (
	"math"

	"github.com/thegreatdaniad/go-tracer/obj_parser"
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
		Vertices: o.Vertices,
	}

	// Copy texture coordinates
	for _, tc := range o.TextureCoordinates {
		newObj.TextureCoordinates = append(newObj.TextureCoordinates, TextureCoordinate{
			U: tc.U,
			V: tc.V,
		})
	}

	// Check if normals are provided
	normalsProvided := len(o.Normals) > 0

	if normalsProvided {
		for _, n := range o.Normals {
			newObj.Normals = append(newObj.Normals, Normal{
				X: float32(n.X),
				Y: float32(n.Y),
				Z: float32(n.Z),
			})
		}
	}

	for _, f := range o.Faces {
		newFace := Face{
			VertexIndices:            make([]int, len(f.VertexIndices)),
			TextureCoordinateIndices: make([]int, len(f.TextureCoordinateIndices)),
			NormalIndices:            make([]int, len(f.VertexIndices)), // Initialize with the length of VertexIndices
		}
		copy(newFace.VertexIndices, f.VertexIndices)
		copy(newFace.TextureCoordinateIndices, f.TextureCoordinateIndices)

		if !normalsProvided {
			// Calculate normal for this face
			normal := calculateFaceNormal(newObj.Vertices, newFace.VertexIndices)
			normalIndex := len(newObj.Normals) // Index of the new normal
			newObj.Normals = append(newObj.Normals, Normal{X: normal[0], Y: normal[1], Z: normal[2]})

			// Assign the normal index to all vertices of the face
			for i := range newFace.NormalIndices {
				newFace.NormalIndices[i] = normalIndex
			}
		} else {
			// If normals are provided, use the original indices
			copy(newFace.NormalIndices, f.NormalIndices)
		}

		newObj.Faces = append(newObj.Faces, newFace)
	}

	return newObj, nil
}
func calculateFaceNormal(vertices []vec3.T, vertexIndices []int) [3]float32 {
	if len(vertexIndices) < 3 {
		return [3]float32{0, 0, 0}
	}

	v0 := vertices[vertexIndices[0]]
	v1 := vertices[vertexIndices[1]]
	v2 := vertices[vertexIndices[2]]

	edge1 := vec3.Sub(&v1, &v0)
	edge2 := vec3.Sub(&v2, &v0)

	normalVec := vec3.Cross(&edge1, &edge2)

	normalVec.Normalize()

	return [3]float32{normalVec[0], normalVec[1], normalVec[2]}
}
func (f Face) Intersects(ray Ray, vertices []vec3.T) (bool, float32) {
	if len(f.VertexIndices) < 3 {
		return false, 0
	}

	v0 := vertices[f.VertexIndices[0]]
	v1 := vertices[f.VertexIndices[1]]
	v2 := vertices[f.VertexIndices[2]]

	// Edge vectors
	e1 := (vec3.Sub(&v1, &v0))
	e2 := vec3.Sub(&v2, &v0)

	// Begin calculating determinant - also used to calculate u parameter
	pvec := vec3.Cross(&ray.Direction, &e2)

	// If determinant is near zero, ray lies in plane of triangle
	det := vec3.Dot(&pvec, &e1)

	// NOT CULLING
	if det > -float32(math.SmallestNonzeroFloat32) && det < float32(math.SmallestNonzeroFloat32) {
		return false, 0
	}
	invDet := 1.0 / det

	// Calculate distance from V0 to ray origin
	tvec := vec3.Sub(&ray.Origin, &v0)

	// Calculate u parameter and test bound
	u := vec3.Dot(&pvec, &tvec) * invDet

	// The intersection lies outside of the triangle
	if u < 0.0 || u > 1.0 {
		return false, 0
	}

	// Prepare to test v parameter
	qvec := vec3.Cross(&tvec, &e1)

	// Calculate V parameter and test bound
	v := vec3.Dot(&qvec, &(ray.Direction)) * invDet

	// The intersection lies outside of the triangle
	if v < 0.0 || u+v > 1.0 {
		return false, 0
	}

	t := vec3.Dot(&qvec, &e2) * invDet

	if t > float32(math.SmallestNonzeroFloat32) {
		return true, t
	}

	return false, 0
}

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
	Vertices []vec3.T
	Faces    []Face
	Material Material
}

type Normal struct {
	X, Y, Z float64
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

	for _, tc := range o.TextureCoordinates {
		newObj.TextureCoordinates = append(newObj.TextureCoordinates, TextureCoordinate{
			U: tc.U,
			V: tc.V,
		})
	}

	for _, n := range o.Normals {
		newObj.Normals = append(newObj.Normals, Normal{
			X: n.X,
			Y: n.Y,
			Z: n.Z,
		})
	}

	for _, f := range o.Faces {
		newFace := Face{
			VertexIndices:            make([]int, len(f.VertexIndices)),
			TextureCoordinateIndices: make([]int, len(f.TextureCoordinateIndices)),
			NormalIndices:            make([]int, len(f.NormalIndices)),
		}
		copy(newFace.VertexIndices, f.VertexIndices)
		copy(newFace.TextureCoordinateIndices, f.TextureCoordinateIndices)
		copy(newFace.NormalIndices, f.NormalIndices)
		newObj.Faces = append(newObj.Faces, newFace)
	}

	return newObj, nil
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

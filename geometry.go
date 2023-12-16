package main

import "github.com/ungerik/go3d/vec3"

type Geometry interface {
	GetGeometryData() GeometryData
}

type GeometryData struct {
	Vertices []vec3.T
	Polygons []Polygon
}
type Polygon struct {
	Vertices [3]vec3.T
}

func (p Polygon) Intersect(r Ray) (bool, float32) {
	v0v1 := vec3.Sub(&p.Vertices[1], &p.Vertices[0])
	v0v2 := vec3.Sub(&p.Vertices[2], &p.Vertices[0])
	normal := vec3.Cross(&v0v1, &v0v2)

	if normal.LengthSqr() < 1e-8 {
		return false, 0
	}
	normal = *normal.Normalize()

	d := -vec3.Dot(&normal, &p.Vertices[0])
	t := -(vec3.Dot(&normal, &r.Origin) + d) / vec3.Dot(&normal, &r.Direction)
	if t < 0 {
		return false, 0 
	}

	P := vec3.Add(&r.Origin, r.Direction.Scale(t))

	edge0 := vec3.Sub(&p.Vertices[1], &p.Vertices[0])
	vp0 := vec3.Sub(&P, &p.Vertices[0])
	C := vec3.Cross(&edge0, &vp0)
	if vec3.Dot(&normal, &C) < 0 {
		return false, 0
	}

	edge1 := vec3.Sub(&p.Vertices[2], &p.Vertices[1])
	vp1 := vec3.Sub(&P, &p.Vertices[1])
	C = vec3.Cross(&edge1, &vp1)
	if vec3.Dot(&normal, &C) < 0 {
		return false, 0
	}

	edge2 := vec3.Sub(&p.Vertices[0], &p.Vertices[2])
	vp2 := vec3.Sub(&P, &p.Vertices[2])
	C = vec3.Cross(&edge2, &vp2)
	if vec3.Dot(&normal, &C) < 0 {
		return false, 0
	}

	// Step 4: Return the intersection distance
	return true, float32(t)
}

type Box struct {
	Width    float32
	Height   float32
	Depth    float32
	Origin   vec3.T
	Material *Material
}

type Sphere struct {
	Radius   float32
	Origin   vec3.T
	Material *Material
}

type Plane struct {
	Origin   vec3.T
	Normal   vec3.T
	Material *Material
}

package main

import (
	"math"

	"github.com/ungerik/go3d/vec3"
)

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
	Vertices []vec3.T
	Polygons []Polygon
	Material *Material
}

func CreateBox(Width float32, Height float32, Depth float32, Origin vec3.T) Box {
	b := Box{
		Width: Width, Height: Height, Depth: Depth, Origin: Origin,
	}
	halfWidth := b.Width / 2
	halfHeight := b.Height / 2
	halfDepth := b.Depth / 2

	vertices := []vec3.T{
		{b.Origin[0] - halfWidth, b.Origin[1] - halfHeight, b.Origin[2] - halfDepth},
		{b.Origin[0] + halfWidth, b.Origin[1] - halfHeight, b.Origin[2] - halfDepth},
		{b.Origin[0] + halfWidth, b.Origin[1] + halfHeight, b.Origin[2] - halfDepth},
		{b.Origin[0] - halfWidth, b.Origin[1] + halfHeight, b.Origin[2] - halfDepth},
		{b.Origin[0] - halfWidth, b.Origin[1] - halfHeight, b.Origin[2] + halfDepth},
		{b.Origin[0] + halfWidth, b.Origin[1] - halfHeight, b.Origin[2] + halfDepth},
		{b.Origin[0] + halfWidth, b.Origin[1] + halfHeight, b.Origin[2] + halfDepth},
		{b.Origin[0] - halfWidth, b.Origin[1] + halfHeight, b.Origin[2] + halfDepth},
	}

	// Step 2: Define the polygons (triangles) for each face of the box
	polygons := []Polygon{
		// Front face
		{Vertices: [3]vec3.T{vertices[0], vertices[1], vertices[2]}},
		{Vertices: [3]vec3.T{vertices[0], vertices[2], vertices[3]}},
		// Back face
		{Vertices: [3]vec3.T{vertices[4], vertices[6], vertices[5]}},
		{Vertices: [3]vec3.T{vertices[4], vertices[7], vertices[6]}},
		// Left face
		{Vertices: [3]vec3.T{vertices[0], vertices[3], vertices[7]}},
		{Vertices: [3]vec3.T{vertices[0], vertices[7], vertices[4]}},
		// Right face
		{Vertices: [3]vec3.T{vertices[1], vertices[5], vertices[6]}},
		{Vertices: [3]vec3.T{vertices[1], vertices[6], vertices[2]}},
		// Top face
		{Vertices: [3]vec3.T{vertices[3], vertices[2], vertices[6]}},
		{Vertices: [3]vec3.T{vertices[3], vertices[6], vertices[7]}},
		// Bottom face
		{Vertices: [3]vec3.T{vertices[0], vertices[5], vertices[1]}},
		{Vertices: [3]vec3.T{vertices[0], vertices[4], vertices[5]}},
	}
	b.Vertices = vertices
	b.Polygons = polygons
	return b
}

type Sphere struct {
	Radius   float32
	Origin   vec3.T
	Material *Material
	Vertices []vec3.T
	Polygons []Polygon
}

func CreateSphere(radius float32, origin vec3.T) Sphere {
	var vertices []vec3.T
	var polygons []Polygon
	resolution := 32
	// Generate vertices
	for lat := 0; lat <= resolution; lat++ {
		theta := float32(lat) * math.Pi / float32(resolution)
		sinTheta := math.Sin(float64(theta))
		cosTheta := math.Cos(float64(theta))

		for lon := 0; lon <= resolution; lon++ {
			phi := float32(lon) * 2 * math.Pi / float32(resolution)
			sinPhi := math.Sin(float64(phi))
			cosPhi := math.Cos(float64(phi))

			x := cosPhi * sinTheta
			y := cosTheta
			z := sinPhi * sinTheta

			vertices = append(vertices, vec3.T{origin[0] + radius*float32(x), origin[1] + radius*float32(y), origin[2] + radius*float32(z)})
		}
	}

	// Generate polygons
	for lat := 0; lat < resolution; lat++ {
		for lon := 0; lon < resolution; lon++ {
			a := lat*(resolution+1) + lon
			b := a + resolution + 1
			polygons = append(polygons, Polygon{Vertices: [3]vec3.T{vertices[a], vertices[b], vertices[a+1]}})
			polygons = append(polygons, Polygon{Vertices: [3]vec3.T{vertices[b], vertices[b+1], vertices[a+1]}})
		}
	}

	return Sphere{
		Radius:   radius,
		Origin:   origin,
		Vertices: vertices,
		Polygons: polygons,
	}
}

type Plane struct {
	Origin   vec3.T
	Normal   vec3.T
	Material *Material
	Vertices []vec3.T
	Polygons []Polygon
}

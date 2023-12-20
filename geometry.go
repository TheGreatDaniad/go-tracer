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

func (p *Polygon) Intersects(ray Ray) (bool, float32) {
	const epsilon = 1e-8

	// Calculate edges
	edge1 := vec3.Sub(&p.Vertices[1], &p.Vertices[0])
	edge2 := vec3.Sub(&p.Vertices[2], &p.Vertices[0])

	// Begin calculating determinant - also used to calculate U parameter
	pvec := vec3.Cross(&ray.Direction, &edge2)

	// If determinant is near zero, ray lies in plane of triangle
	det := vec3.Dot(&edge1, &pvec)

	// NOT CULLING
	if math.Abs(float64(det)) < epsilon {
		return false, 0
	}

	invDet := 1.0 / det

	// Calculate distance from V1 to ray origin
	tvec := vec3.Sub(&ray.Origin, &p.Vertices[0])

	// Calculate U parameter and test bound
	u := vec3.Dot(&tvec, &pvec) * invDet
	if u < 0.0 || u > 1.0 {
		return false, 0
	}

	// Prepare to test V parameter
	qvec := vec3.Cross(&tvec, &edge1)

	// Calculate V parameter and test bound
	v := vec3.Dot(&ray.Direction, &qvec) * invDet
	if v < 0.0 || u+v > 1.0 {
		return false, 0
	}

	// Calculate t, the distance from the ray origin to the intersection point
	t := vec3.Dot(&edge2, &qvec) * invDet

	return true, t
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
func (b Box) GetGeometryData() GeometryData {
	return GeometryData{Vertices: b.Vertices, Polygons: b.Polygons}
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

func (s Sphere) GetGeometryData() GeometryData {
	return GeometryData{Vertices: s.Vertices, Polygons: s.Polygons}
}

type Plane struct {
	Origin   vec3.T
	Normal   vec3.T
	Material *Material
	Vertices []vec3.T
	Polygons []Polygon
}

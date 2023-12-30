package main

import (
	"image/color"

	"github.com/ungerik/go3d/vec3"
)

type Light struct {
	Position    vec3.T
	Color       color.RGBA
	Intensity   float32
	Attenuation float32
}

func CreateLight(position vec3.T, color color.RGBA, intensity, attenuation float32) Light {
	return Light{position, color, intensity, attenuation}
}

func (l Light) CalculateColorContribution(geo Geometry, point vec3.T, face Face) color.RGBA {
	obj := geo.GetGeometryData()
	lightDir := vec3.Sub(&l.Position, &point)
	lightDir.Normalize()

	// Retrieve the vertices of the face
	V1 := obj.Vertices[face.VertexIndices[0]]
	V2 := obj.Vertices[face.VertexIndices[1]]
	V3 := obj.Vertices[face.VertexIndices[2]]

	// Compute vectors for two edges of the triangle
	edge1 := vec3.Sub(&V2, &V1)
	edge2 := vec3.Sub(&V3, &V1)
	vp := vec3.Sub(&point, &V1)

	// Compute dot products
	dot00 := vec3.Dot(&edge1, &edge1)
	dot01 := vec3.Dot(&edge1, &edge2)
	dot02 := vec3.Dot(&edge1, &vp)
	dot11 := vec3.Dot(&edge2, &edge2)
	dot12 := vec3.Dot(&edge2, &vp)

	// Compute barycentric coordinates
	invDenom := 1 / (dot00*dot11 - dot01*dot01)
	u := (dot11*dot02 - dot01*dot12) * invDenom
	v := (dot00*dot12 - dot01*dot02) * invDenom
	w := 1.0 - u - v

	// Interpolate normal using the barycentric coordinates
	interpolatedNormal := vec3.T{
		u*obj.Normals[face.NormalIndices[0]].X + v*obj.Normals[face.NormalIndices[1]].X + w*obj.Normals[face.NormalIndices[2]].X,
		u*obj.Normals[face.NormalIndices[0]].Y + v*obj.Normals[face.NormalIndices[1]].Y + w*obj.Normals[face.NormalIndices[2]].Y,
		u*obj.Normals[face.NormalIndices[0]].Z + v*obj.Normals[face.NormalIndices[1]].Z + w*obj.Normals[face.NormalIndices[2]].Z,
	}
	interpolatedNormal.Normalize()
	cosTheta := vec3.Dot(&interpolatedNormal, &lightDir)
	if cosTheta < 0 {
		cosTheta = 0
	}

	// Apply the light intensity and attenuation with the cosine of the angle
	// between the light direction and the interpolated normal
	r := float32(l.Color.R) * cosTheta * l.Intensity / l.Attenuation
	g := float32(l.Color.G) * cosTheta * l.Intensity / l.Attenuation
	b := float32(l.Color.B) * cosTheta * l.Intensity / l.Attenuation
	a := float32(l.Color.A)

	// Clamp the color components to the range [0, 255]
	r = clampColorComponent(r)
	g = clampColorComponent(g)
	b = clampColorComponent(b)

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}

	// Continue with the shading calculation using interpolated normal
	// ...
}

// Helper function to clamp color component to 0-255 range
func clampColorComponent(value float32) float32 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return value
}

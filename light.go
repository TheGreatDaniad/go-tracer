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

	u, v, w := ComputeBarycentricCoordinates(point, V1, V2, V3)

	// Interpolate the normal using the barycentric coordinates
	interpolatedNormal := vec3.T{
		u*obj.Normals[face.NormalIndices[0]].X + v*obj.Normals[face.NormalIndices[1]].X + w*obj.Normals[face.NormalIndices[2]].X,
		u*obj.Normals[face.NormalIndices[0]].Y + v*obj.Normals[face.NormalIndices[1]].Y + w*obj.Normals[face.NormalIndices[2]].Y,
		u*obj.Normals[face.NormalIndices[0]].Z + v*obj.Normals[face.NormalIndices[1]].Z + w*obj.Normals[face.NormalIndices[2]].Z,
	}
	interpolatedNormal.Normalize()

	// Dot product to find the cosine of the angle between the light and the interpolated normal
	cosTheta := vec3.Dot(&interpolatedNormal, &lightDir)
	if cosTheta < 0 {
		cosTheta = 0
	}

	r := float32(l.Color.R) * cosTheta * l.Intensity / l.Attenuation
	g := float32(l.Color.G) * cosTheta * l.Intensity / l.Attenuation
	b := float32(l.Color.B) * cosTheta * l.Intensity / l.Attenuation
	a := float32(l.Color.A)

	// Clamp the color components to the range [0, 255]
	r = clampColorComponent(r)
	g = clampColorComponent(g)
	b = clampColorComponent(b)

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

func clampColorComponent(value float32) float32 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return value
}
func ComputeBarycentricCoordinates(P, A, B, C vec3.T) (u, v, w float32) {
	// Vectors from A to B and A to C
	v0 := vec3.Sub(&B, &A)
	v1 := vec3.Sub(&C, &A)
	v2 := vec3.Sub(&P, &A)

	// Compute dot products
	d00 := vec3.Dot(&v0, &v0)
	d01 := vec3.Dot(&v0, &v1)
	d11 := vec3.Dot(&v1, &v1)
	d20 := vec3.Dot(&v2, &v0)
	d21 := vec3.Dot(&v2, &v1)

	// Compute denominator
	denom := d00*d11 - d01*d01
	v = (d11*d20 - d01*d21) / denom
	w = (d00*d21 - d01*d20) / denom
	u = 1.0 - v - w

	return u, v, w
}

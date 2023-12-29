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
	// Calculate the normalized light direction
	lightDir := vec3.Sub(&l.Position, &point)
	lightDir.Normalize()

	// Calculate the normal of the face
	// Assuming the face normal is the average of the normals of the vertices
	var faceNormal vec3.T
	for _, ni := range face.NormalIndices {
		temp := obj.Normals[ni].ToVec3()
		faceNormal.Add(&temp)
	}
	faceNormal.Scale(1.0 / float32(len(face.NormalIndices)))
	faceNormal.Normalize()

	// Dot product to find the cosine of the angle between the light and the face normal
	cosTheta := vec3.Dot(&faceNormal, &lightDir)
	// Clamp the value between 0 and 1
	if cosTheta < 0 {
		cosTheta = 0
	}

	r := float32(l.Color.R) * cosTheta * l.Intensity / l.Attenuation
	g := float32(l.Color.G) * cosTheta * l.Intensity / l.Attenuation
	b := float32(l.Color.B) * cosTheta * l.Intensity / l.Attenuation
	a := float32(l.Color.A)

	r = clampColorComponent(r)
	g = clampColorComponent(g)
	b = clampColorComponent(b)
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
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

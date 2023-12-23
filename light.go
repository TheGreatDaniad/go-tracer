package main

import (
	"fmt"
	"image/color"

	"github.com/ungerik/go3d/vec3"
)

type Light struct {
	Position    vec3.T
	Color       color.RGBA
	Intensity   float32
	Attenuation float32
}

// creates a basic light
func CreateLight(position vec3.T, color color.RGBA, intensity, attenuation float32) Light {
	return Light{position, color, intensity, attenuation}
}

func (l Light) CalculateColorContribution(point, normal vec3.T, surfaceColor color.RGBA) color.RGBA {
	// Direction from point to light
	lightDir := vec3.Sub(&l.Position, &point)
	lightDir = *lightDir.Normalize()

	// Dot product of light direction and normal - Lambert's Cosine Law
	dot := vec3.Dot(&normal, &lightDir)
	if dot > 0 {
		diffuseIntensity := dot * l.Intensity * 10
		// Attenuation (optional, can be omitted or modified as needed)
		distance := vec3.Distance(&l.Position, &point)
		attenuation := 1 / (1 + l.Attenuation*distance*distance)
		diffuseIntensity *= attenuation

		// Calculate the contribution of the light source
		lightContribution := vec3.T{
			float32(surfaceColor.R*l.Color.R) * diffuseIntensity,
			float32(surfaceColor.G*l.Color.G) * diffuseIntensity,
			float32(surfaceColor.B*l.Color.B) * diffuseIntensity,
		}
		if lightContribution[0] > 255 {
			lightContribution[0] = 255
		}
		if lightContribution[1] > 255 {
			lightContribution[1] = 255
		}
		if lightContribution[2] > 255 {
			lightContribution[2] = 255
		}
		color := color.RGBA{
			R: uint8(lightContribution[0]),
			G: uint8(lightContribution[1]),
			B: uint8(lightContribution[2]),
		}
		fmt.Println(color)
		return color
	}

	return color.RGBA{0, 0, 0, 0}
}

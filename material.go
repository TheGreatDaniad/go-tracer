package main

import "image/color"

type Material struct {
	Color        color.RGBA
	Reflectivity float32
	Opacity      float32
	Diffuse      float32
	Roughness    float32
}

package main

import "github.com/ungerik/go3d/vec3"

type Material struct {
	Color        vec3.T
	Reflectivity float32
	Opacity      float32
	Diffuse      float32
	Roughness    float32
}

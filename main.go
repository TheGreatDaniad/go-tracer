package main

import "github.com/ungerik/go3d/vec3"

func main() {
	s1 := Sphere{Radius: 1, Origin: vec3.T{0, 0, 0}, Material: &Material{Color: vec3.T{1, 0, 0}, Reflectivity: 0.5, Opacity: 0.5, Diffuse: 0.5, Roughness: 0.5}}
	cam := CreateCamera(vec3.T{-4, 0, 0}, vec3.T{1, 0, 0}, 90, 1)
	
}

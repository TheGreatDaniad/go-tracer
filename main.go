package main

import "github.com/ungerik/go3d/vec3"

func main() {
	s1 := CreateSphere(1, vec3.T{0, 0, 0})
	cam := CreateCamera(vec3.T{-4, 0, 0}, vec3.T{1, 0, 0}, 90, 1, 1024, 768)
	
}

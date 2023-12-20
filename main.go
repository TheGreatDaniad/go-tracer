package main

import "github.com/ungerik/go3d/vec3"

func main() {
	var s Space

	camera := CreateCamera(vec3.T{0, 0, 0}, vec3.T{1, 0, 0}, 90, 1, 100, 100)
	sphere := CreateSphere(1, vec3.T{0, 0, 1})
	s.AddGeometry(&sphere)
	camera.Render(&s)

}

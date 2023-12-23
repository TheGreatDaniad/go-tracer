package main

import (
	"image/color"
	"time"

	"github.com/ungerik/go3d/vec3"
)

func main() {
	var s Space
	now := time.Now()
	camera := CreateCamera(vec3.T{0, 0, 0}, vec3.T{1, 0, 0}, vec3.T{0, 1, 0}, 90, 1, 600, 600)
	light := CreateLight(vec3.T{-3, 5, 0}, color.RGBA{255, 255, 255, 255}, 1, 0.1)
	s.AddLight(light)
	sphere := CreateSphere(2, vec3.T{5, 0, 0})
	sphere.SetMaterial(Material{color.RGBA{100, 100, 100, 255}, 0.5, 0, 0.5, 0.5})

	s.AddGeometry(&sphere)
	camera.Render(&s)
	elapsed := time.Since(now)
	println(elapsed.Seconds())

}

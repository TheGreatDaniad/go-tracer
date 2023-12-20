package main

import (
	"time"

	"github.com/ungerik/go3d/vec3"
)

func main() {
	var s Space
	now := time.Now()
	camera := CreateCamera(vec3.T{0, 0, 0}, vec3.T{1, 0, 0}, vec3.T{0, 1, 0}, 90, 1, 600, 600)
	sphere := CreateSphere(1, vec3.T{5, 0, 0})
	s.AddGeometry(&sphere)
	camera.Render(&s)
	elapsed := time.Since(now)
	println(elapsed.Seconds())

}

package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/ungerik/go3d/vec3"
)

// func main() {
// 	var s Space
// 	now := time.Now()
// 	camera := CreateCamera(vec3.T{0, 0, 0}, vec3.T{1, 0, 0}, vec3.T{0, 1, 0}, 90, 1, 600, 600)
// 	light := CreateLight(vec3.T{-3, 5, 0}, color.RGBA{255, 255, 255, 255}, 1, 0.1)
// 	s.AddLight(light)
// 	sphere := CreateSphere(1, vec3.T{5, 0, 0})
// 	sphere.SetMaterial(Material{color.RGBA{100, 100, 100, 255}, 0.5, 0, 0.5, 0.5})
// 	// box := CreateBox(2, 2, 2, vec3.T{5, 2, 3})
// 	// box.SetMaterial(Material{color.RGBA{100, 100, 100, 255}, 0.5, 0, 0.5, 0.5})
// 	// s.AddGeometry(&box)

// 	s.AddGeometry(&sphere)
// 	camera.Render(&s)
// 	elapsed := time.Since(now)
// 	println(elapsed.Seconds())

// }
func main() {
	var s Space
	now := time.Now()
	camera := CreateCamera(vec3.T{-5, 2, 1}, vec3.T{1, 0.0, 0.0}, vec3.T{0, -1, 0}, 90, 1, 200, 200)
	light := CreateLight(vec3.T{-3, 7, 2}, color.RGBA{214, 153, 88, 255}, 0.2, 0.2)
	s.AddLight(light)
	o, err := ParseObjFile("objects/meerschaum_new.obj")
	if err != nil {
		panic(err)
	}
	o.Rotate(-20,90,0)
	s.AddGeometry(o)
	camera.Render(&s)
	fmt.Println(time.Since(now).Seconds())

}

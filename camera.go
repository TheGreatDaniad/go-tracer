package main

import (
	"math"

	"github.com/ungerik/go3d/vec3"
)

type Camera struct {
	Origin                vec3.T
	Direction             vec3.T
	FieldOfViewHorizontal float32
	AspectRatio           float32
	Width                 float32
	Height                float32
	FocalLength           float32
}

func (c Camera) CalculateFocalLength(sensorWidth, fov float64) float64 {
	return (sensorWidth / 2) / math.Tan(fov/2*(math.Pi/180))
}
func CreateCamera(origin, direction vec3.T, fov, aspectRatio float32) Camera {
	c := Camera{}
	c.Origin = origin
	c.Direction = direction
	c.FieldOfViewHorizontal = fov
	c.AspectRatio = aspectRatio
	c.Width = 1
	c.Height = c.Width / c.AspectRatio
	c.FocalLength = float32(c.CalculateFocalLength(float64(c.Width), float64(c.FieldOfViewHorizontal)))
	return c
}


type Ray struct {	
	Origin vec3.T
	Direction vec3.T
	 	
}

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
	ResolutionX           int
	ResolutionY           int
}

func (c Camera) CalculateFocalLength(sensorWidth, fov float64) float64 {
	return (sensorWidth / 2) / math.Tan(fov/2*(math.Pi/180))
}
func CreateCamera(origin, direction vec3.T, fov, aspectRatio float32, resolutionX, resolutionY int) Camera {
	c := Camera{}
	c.Origin = origin
	c.Direction = direction
	c.FieldOfViewHorizontal = fov
	c.AspectRatio = aspectRatio
	c.Width = 1
	c.Height = c.Width / c.AspectRatio
	c.FocalLength = float32(c.CalculateFocalLength(float64(c.Width), float64(c.FieldOfViewHorizontal)))
	c.ResolutionX = resolutionX
	c.ResolutionY = resolutionY
	return c
}

func (camera Camera) CreateRays() []Ray {
	rays := make([]Ray, camera.ResolutionY)
	for y := 0; y < camera.ResolutionY; y++ {
		for x := 0; x < camera.ResolutionX; x++ {
			direction := calculateRayDirection(camera, x, y)
			rays = append(rays, Ray{Origin: camera.Origin, Direction: direction})
		}
	}
	return rays
}

func calculateRayDirection(camera Camera, pixelX, pixelY int) vec3.T {
	ndcX := (float32(pixelX) + 0.5) / float32(camera.ResolutionX)
	ndcY := (float32(pixelY) + 0.5) / float32(camera.ResolutionY)
	screenX := 2.0*ndcX - 1.0
	screenY := 1.0 - 2.0*ndcY
	screenX *= camera.AspectRatio * camera.Width / camera.Height

	direction := vec3.T{screenX, screenY, -camera.FocalLength}
	direction = direction.Normalized()

	crossResult := vec3.Cross(&vec3.T{0, 0, -1}, &camera.Direction)
	rotationAxis := crossResult.Normalize()
	rotationAngle := math.Acos(float64(vec3.Dot(&vec3.T{0, 0, -1}, camera.Direction.Normalize())))
	direction = rotateVector(direction, *rotationAxis, float32(rotationAngle))

	return direction
}

func rotateVector(v, axis vec3.T, angle float32) vec3.T {
	cosAngle := math.Cos(float64(angle))
	sinAngle := math.Sin(float64(angle))

	term1 := v.Scaled(float32(cosAngle))
	crossProduct := vec3.Cross(&axis, &v)          
	term2 := crossProduct.Scaled(float32(sinAngle)) 

	dotProduct := vec3.Dot(&axis, &v)
	term3 := axis.Scaled(dotProduct * (1 - float32(cosAngle)))

	return *term1.Add(&term2).Add(&term3)
}

type Ray struct {
	Origin    vec3.T
	Direction vec3.T
}

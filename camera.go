package main

import (
	"image"
	"image/color"
	"math"

	"github.com/ungerik/go3d/vec3"
)

type Camera struct {
	Origin                vec3.T
	Direction             vec3.T
	Up                    vec3.T
	FieldOfViewHorizontal float32
	AspectRatio           float32
	Width                 float32
	Height                float32
	FocalLength           float32
	ResolutionX           int
	ResolutionY           int
}

func AddColors(c1, c2 color.RGBA) color.RGBA {
	R := c1.R + c2.R
	if R > 255 {
		R = 255
	}
	G := c1.G + c2.G
	if G > 255 {
		G = 255
	}
	B := c1.B + c2.B
	if B > 255 {
		B = 255
	}
	return color.RGBA{R, G, B, 255}
}
func (c *Camera) CalculatePixelPosition(x, y int) vec3.T {
	aspectRatio := float32(c.ResolutionX) / float32(c.ResolutionY)
	pixelSizeX := (2.0 * c.FocalLength * aspectRatio) / float32(c.ResolutionX)
	pixelSizeY := (2.0 * c.FocalLength) / float32(c.ResolutionY)

	pixelPosX := (float32(x) - float32(c.ResolutionX)/2.0) * pixelSizeX
	pixelPosY := (float32(y) - float32(c.ResolutionY)/2.0) * pixelSizeY

	// Assuming vec3 library provides these methods
	right := vec3.Cross(&c.Direction, &c.Up)
	right = *right.Normalize()

	// Temporary variables for intermediate results
	focalPoint := c.Direction.Scaled(c.FocalLength)
	rightOffset := right.Scaled(pixelPosX)
	upOffset := c.Up.Scaled(pixelPosY)

	// Calculate the pixel world position
	pixelWorldPos := c.Origin.Added(&focalPoint)
	r := pixelWorldPos.Add(&rightOffset)
	r = r.Add(&upOffset)
	return *r
}

type RayFaceIntersection struct {
	ReflectionRay        Ray
	IntersectionPoint    vec3.T
	IntersectionDistance float32
	Face                 Face
	Geometry             Geometry
	Material             Material
}

func (c Camera) CalculateFocalLength(sensorWidth, fov float64) float64 {
	return (sensorWidth / 2) / math.Tan(fov/2*(math.Pi/180))
}

func (c Camera) Render(s *Space) {
	img := image.NewRGBA(image.Rect(0, 0, c.ResolutionX, c.ResolutionY))
	rays := c.CreateRays()
	for i, ray := range rays {
		var intersections = []RayFaceIntersection{}
		for _, geometry := range s.Geometries {
			Faces := (*geometry).GetGeometryData().Faces
			for _, Face := range Faces {
				intersects, dis, point := Face.Intersects(ray, (*geometry).GetGeometryData().Vertices)
				if intersects {
					intersection := RayFaceIntersection{IntersectionPoint: point, IntersectionDistance: dis, Face: Face, Material: ((*geometry).GetGeometryData().Material), Geometry: *geometry}
					intersections = append(intersections, intersection)
				}
			}

		}
		if len(intersections) > 0 {
			var finalColor color.RGBA
			intersection := intersections[0]
			for _, i := range intersections {
				if i.IntersectionDistance < intersection.IntersectionDistance {
					intersection = i
				}
			}
			for _, light := range s.Lights {
				lightContribution := light.CalculateColorContribution(intersection.Geometry, intersection.IntersectionPoint, intersection.Face)
				finalColor = AddColors(finalColor, lightContribution)
			}
			finalColor = AddColors(finalColor, intersection.Material.Color)
			img.Set(i%c.ResolutionX, i/c.ResolutionY, finalColor)
		}
	}
	SaveImage(img, "test.png")
}

func CreateCamera(origin, direction vec3.T, up vec3.T, fov, aspectRatio float32, resolutionX, resolutionY int) Camera {
	c := Camera{}
	c.Origin = origin
	c.Direction = direction
	c.Up = up
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
	pixelPosition := camera.CalculatePixelPosition(pixelX, pixelY)
	direction := vec3.Sub(&pixelPosition, &camera.Origin)
	return direction
}

type Ray struct {
	Origin    vec3.T
	Direction vec3.T
}

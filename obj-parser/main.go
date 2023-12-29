package obj_parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ungerik/go3d/vec3"
)

type Normal struct {
	X, Y, Z float64
}

type TextureCoordinate struct {
	U, V float64
}

type Face struct {
	VertexIndices            []int
	TextureCoordinateIndices []int
	NormalIndices            []int
}

type Obj struct {
	Vertices           []vec3.T
	TextureCoordinates []TextureCoordinate
	Normals            []Normal
	Faces              []Face
}

func ParseObjFile(filename string) (*Obj, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	obj := Obj{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		fields := strings.Fields(line)
		switch fields[0] {
		case "v": // Vertex
			vertex, err := parseVertex(fields)
			if err != nil {
				return nil, err
			}
			obj.Vertices = append(obj.Vertices, vertex)
		case "vt": // Texture coordinate
			tc, err := parseTextureCoordinate(fields)
			if err != nil {
				return nil, err
			}
			obj.TextureCoordinates = append(obj.TextureCoordinates, tc)
		case "vn": // Normal
			normal, err := parseNormal(fields)
			if err != nil {
				return nil, err
			}
			obj.Normals = append(obj.Normals, normal)
		case "f": // Face
			face, err := parseFace(fields)
			if err != nil {
				return nil, err
			}
			obj.Faces = append(obj.Faces, face)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &obj, nil
}

func parseVertex(fields []string) (vec3.T, error) {
	if len(fields) < 4 {
		return vec3.T{}, fmt.Errorf("invalid vertex definition: %v", fields)
	}
	x, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return vec3.T{}, err
	}
	y, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return vec3.T{}, err
	}
	z, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return vec3.T{}, err
	}
	return vec3.T{float32(x), float32(y), float32(z)}, nil
}

func parseTextureCoordinate(fields []string) (TextureCoordinate, error) {
	if len(fields) < 3 {
		return TextureCoordinate{}, fmt.Errorf("invalid texture coordinate definition: %v", fields)
	}
	u, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return TextureCoordinate{}, err
	}
	v, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return TextureCoordinate{}, err
	}
	return TextureCoordinate{U: u, V: v}, nil
}

func parseNormal(fields []string) (Normal, error) {
	if len(fields) < 4 {
		return Normal{}, fmt.Errorf("invalid normal definition: %v", fields)
	}
	x, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return Normal{}, err
	}
	y, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return Normal{}, err
	}
	z, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return Normal{}, err
	}
	return Normal{X: x, Y: y, Z: z}, nil
}

func parseFace(fields []string) (Face, error) {
	face := Face{
		VertexIndices:            make([]int, 0, len(fields)-1),
		TextureCoordinateIndices: make([]int, 0, len(fields)-1),
		NormalIndices:            make([]int, 0, len(fields)-1),
	}

	for _, f := range fields[1:] {
		parts := strings.Split(f, "/")
		vi, err := strconv.Atoi(parts[0])
		if err != nil {
			return Face{}, err
		}
		face.VertexIndices = append(face.VertexIndices, vi-1) // OBJ indices are 1-based

		if len(parts) > 1 && parts[1] != "" {
			ti, err := strconv.Atoi(parts[1])
			if err != nil {
				return Face{}, err
			}
			face.TextureCoordinateIndices = append(face.TextureCoordinateIndices, ti-1)
		}

		if len(parts) > 2 && parts[2] != "" {
			ni, err := strconv.Atoi(parts[2])
			if err != nil {
				return Face{}, err
			}
			face.NormalIndices = append(face.NormalIndices, ni-1)
		}
	}

	return face, nil
}

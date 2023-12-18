package main

type Space struct {
	Geometries []*Geometry
}

func (s *Space) AddGeometry(g *Geometry) {
	s.Geometries = append(s.Geometries, g)
}

package main

type Space struct {
	Geometries []*Geometry
	Lights     []*Light
}

func (s *Space) AddGeometry(g Geometry) {
	s.Geometries = append(s.Geometries, &g)
}

func (s *Space) AddLight(l Light) {
	s.Lights = append(s.Lights, &l)
}

package health

import "testing"

type t1 struct {
	h bool
}

func (t *t1) HealthComponentName() string {
	return "t1"
}

func (t *t1) CheckHealthy() (bool, string) {
	return t.h, "foo"
}

func TestSimple(t *testing.T) {
	c1 := &t1{h: true}
	AddComponent(c1)

	s1 := GatherStatus()
	if !s1.Healthy {
		t.Fatal()
	}

	c2 := &t1{h: false}
	AddComponent(c2)

	s2 := GatherStatus()
	if s2.Healthy {
		t.Fatal()
	}
}

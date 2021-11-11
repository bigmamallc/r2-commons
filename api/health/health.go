package health

import "sync"

type Component interface {
	HealthComponentName() string
	CheckHealthy() (bool, string)
}

type tracker struct {
	lock sync.RWMutex

	items []Component
}

var t = &tracker{}

func AddComponent(i Component) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.items = append(t.items, i)
}

type StatusComponent struct {
	Name    string `json:"name"`
	Healthy bool   `json:"healthy"`
	Status  string `json:"status"`
}

type Status struct {
	Healthy    bool               `json:"healthy"`
	Components []*StatusComponent `json:"components"`
}

func GatherStatus() *Status {
	var items []Component
	func() {
		t.lock.RLock()
		defer t.lock.RUnlock()

		items = append(items, t.items...)
	}()

	s := &Status{
		Healthy: true,
	}
	for _, i := range items {
		healthy, status := i.CheckHealthy()
		s.Components = append(s.Components, &StatusComponent{Name: i.HealthComponentName(), Healthy: healthy, Status: status})
		if !healthy {
			s.Healthy = false
		}
	}

	return s
}

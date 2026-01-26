package core

import (
	"sort"
	"sync"
	"time"
)

type GlobalState struct {
	mu       sync.RWMutex
	Profiles map[string]Profile
	States   map[string]State
}

func NewState() *GlobalState {
	return &GlobalState{
		Profiles: make(map[string]Profile),
		States:   make(map[string]State),
	}
}

func (s *GlobalState) UpdateFromPayload(payload InitialStatePayload) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range payload.Profiles {
		s.Profiles[p.ID] = p
	}

	for _, st := range payload.States {
		s.States[st.ID] = st
	}
}

type ViewItem struct {
	Name       string
	Battery    float64
	IsCharging bool
	Lat        float64
	Lon        float64
	LastSeen   time.Time
}

func (s *GlobalState) GetViewData() []ViewItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var data []ViewItem

	for id, state := range s.States {
		name := "Unknown " + id
		if _, ok := s.Profiles[id]; ok {
			name = s.Profiles[id].Name
		}

		item := ViewItem{
			Name:       name,
			Battery:    state.Battery.Level * 100,
			IsCharging: state.Battery.IsCharging,
			Lat:        state.Location.Lat,
			Lon:        state.Location.Lon,
			LastSeen:   time.Unix(state.LastSeen/1000, 0),
		}

		data = append(data, item)
	}

	// sorting is used to keep order of friends the same on each render
	sort.Slice(data, func(i, j int) bool {
		return data[i].Name < data[j].Name
	})

	return data
}

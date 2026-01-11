package state

import (
	"sort"
	"sync"
	"time"

	"github.com/rangaroo/2gis-friends/internal/models"
)

type GlobalStore struct {
	mu          sync.RWMutex
	Profiles    map[string]models.Profile
	States      map[string]models.State
	IsConnected bool
}

func NewStore() *GlobalStore {
	return &GlobalStore{
		Profiles: make(map[string]models.Profile),
		States:   make(map[string]models.State),
	}
}

func (s *GlobalStore) UpdateFromPayload(payload models.InitialStatePayload) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range payload.Profiles {
		s.Profiles[p.ID] = p
	}

	for _, st := range payload.States {
		s.States[st.ID] = st
	}
}

func (s *GlobalStore) SetConnection(connected bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IsConnected = connected
}

type ViewItem struct {
	Name       string
	Battery    float64
	IsCharging bool
	Lat        float64
	Lon        float64
	LastSeen   time.Time
}

func (s *GlobalStore) GetViewData() []ViewItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var data []ViewItem

	for id, state := range s.States {
		name := "Unknown " + id
		if _, ok := s.Profiles[id]; ok {
			name = s.Profiles[id].Name
		}

		lastSeen := time.Unix(state.LastSeen/1000, 0)

		data = append(data, ViewItem{
			Name:       name,
			Battery:    state.Battery.Level * 100,
			IsCharging: state.Battery.IsCharging,
			Lat:        state.Location.Lat,
			Lon:        state.Location.Lon,
			LastSeen:   lastSeen,
		})
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Name < data[j].Name
	})

	return data
}

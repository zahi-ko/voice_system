package memory

import (
	"sync"
	"voice_system/internal/domain/effects"
)

var effect_store = map[string]func() effects.Effect{
    "gain":      func() effects.Effect { return &effects.EffectGain{} },
    "tempo":     func() effects.Effect { return &effects.EffectTempo{} },
    "reverse":   func() effects.Effect { return &effects.EffectReverse{} },
    "tempo_pv":  func() effects.Effect { return &effects.EffectTempoPV{} },
    "normalize": func() effects.Effect { return &effects.EffectNormalization{} },
}

type Registry struct {
	mu          sync.RWMutex
	effectStore map[string]func() effects.Effect
}

func NewRegistry() (*Registry, error) {
	return &Registry{
		effectStore: effect_store,
	}, nil
}

func (r *Registry) Register(name string, effect effects.Effect) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.effectStore[name] = func() effects.Effect { return effect }
}

func (r *Registry) Get(name string) (effects.Effect, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	effect, ok := r.effectStore[name]
	return effect(), ok
}

func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.effectStore))

	for name := range r.effectStore {
		names = append(names, name)
	}

	return names
}

func (r *Registry) GetMap() map[string]func() effects.Effect {
	return r.effectStore
}

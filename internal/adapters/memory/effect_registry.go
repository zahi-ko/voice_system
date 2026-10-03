package memory

import (
	"fmt"
	"sync"
	"voice_system/internal/domain/effects"
)

var effect_store = map[string]func(effects.Payload) (effects.Effect, error){
	"gain":      effects.NewGainFromPayload,
	"tempo":     effects.NewTempoFromPayload,
	"reverse":   effects.NewReverseFromPayload,
	"tempo_pv":  effects.NewTempoPVFromPayload,
	"normalize": effects.NewNormalizationFromPayload,
}

type Registry struct {
	mu          sync.RWMutex
	effectStore map[string]func(effects.Payload) (effects.Effect, error)
}

func NewRegistry() (*Registry, error) {
	return &Registry{
		effectStore: effect_store,
	}, nil
}

// Bind 按名字构造效果，各效果的 FromPayload 构造函数负责参数校验。
func (r *Registry) Bind(payload effects.Payload) (effects.Effect, error) {
	rawName, ok := payload["name"]
	if !ok {
		return nil, fmt.Errorf("effect name not provided in payload")
	}

	name, ok := rawName.(string)
	if !ok {
		return nil, fmt.Errorf("effect name must be a string")
	}

	r.mu.RLock()
	factory, ok := r.effectStore[name]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("effect %s not found", name)
	}

	return factory(payload)
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

package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"voice_system/internal/application/ports"
	"voice_system/internal/domain/audio"
	"voice_system/internal/domain/effects"

	"github.com/google/uuid"
)

type EffectService struct {
	last     string
	store    ports.AudioStore
	registry ports.EffectRegistry
}

func NewEffectService(
	store ports.AudioStore,
	registry ports.EffectRegistry,
) *EffectService {
	return &EffectService{
		store:    store,
		registry: registry,
	}
}

func (s *EffectService) Apply(ctx context.Context, audioId string, save bool, payload effects.Payload) (audio.Audio, error) {
	effect, err := s.bindEffect(payload)
	if err != nil {
		return audio.Audio{}, fmt.Errorf("bind effect: %w", err)
	}

	aud, data, err := s.store.Get(ctx, audioId)
	if err != nil {
		return audio.Audio{}, fmt.Errorf("get original audio: %w", err)
	}

	out, meta, err := effect.Apply(data, aud.Meta)
	if err != nil {
		return audio.Audio{}, fmt.Errorf("apply effect %s: %w", effect.Name(), err)
	}

	newAud := audio.Audio{
		ID:   audioId,
		Meta: meta,
		Name: aud.Name + "_" + effect.Name(),
	}

	if !save {
		if err := s.store.Replace(ctx, audioId, newAud, out); err != nil {
			return audio.Audio{}, fmt.Errorf("replace audio: %w", err)
		}
	} else {
		if err := s.store.Save(ctx, newAud, out); err != nil {
			return audio.Audio{}, fmt.Errorf("save new audio: %w", err)
		}
		newAud.ID = uuid.NewString()
	}

	return newAud, nil
}

// ApplyChain 按顺序应用效果链，每一步以前一步输出为输入。
func (s *EffectService) ApplyChain(ctx context.Context, audioId string, save bool, payloads []effects.Payload) (audio.Audio, string, error) {
	applied := ""

	aud, data, err := s.store.Get(ctx, audioId)
	if err != nil {
		return audio.Audio{}, "", fmt.Errorf("get original audio: %w", err)
	}

	for _, payload := range payloads {
		effect, err := s.bindEffect(payload)
		if err != nil {
			return audio.Audio{}, "", fmt.Errorf("bind effect: %w", err)
		}

		data, aud.Meta, err = effect.Apply(data, aud.Meta)
		if err != nil {
			return audio.Audio{}, "", fmt.Errorf("apply effect %s: %w", effect.Name(), err)
		}

		applied += effect.Name() + "->"
	}

	newAud := audio.Audio{
		ID:   audioId,
		Meta: aud.Meta,
		Name: aud.Name + "_" + applied,
	}

	if !save {
		if err := s.store.Replace(ctx, audioId, newAud, data); err != nil {
			return audio.Audio{}, "", fmt.Errorf("replace audio: %w", err)
		}
	} else {
		if err := s.store.Save(ctx, newAud, data); err != nil {
			return audio.Audio{}, "", fmt.Errorf("save new audio: %w", err)
		}
		newAud.ID = uuid.NewString()
	}

	return newAud, applied, nil

}

func (s *EffectService) ListAvailable() []string {
	return s.registry.Names()
}

func (s *EffectService) bindEffect(payload effects.Payload) (effects.Effect, error) {
	effectMap := s.registry.GetMap()

	rawName, ok := payload["name"]
	if !ok {
		return nil, errors.New("effect name not provided in payload")
	}

	name, ok := rawName.(string)
	if !ok {
		return nil, errors.New("effect name must be a string")
	}

	eff, ok := effectMap[name]
	if !ok {
		return nil, fmt.Errorf("effect %s not found", name)
	}

	effInstance := eff()
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	if err := json.Unmarshal(data, effInstance); err != nil {
		return nil, fmt.Errorf("unmarshal payload to effect: %w", err)
	}

	return effInstance, nil
}

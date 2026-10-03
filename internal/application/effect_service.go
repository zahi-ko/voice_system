package application

import (
	"context"
	"fmt"

	"voice_system/internal/application/ports"
	"voice_system/internal/domain/audio"
	"voice_system/internal/domain/effects"

	"github.com/google/uuid"
)

type EffectService struct {
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
	effect, err := s.registry.Bind(payload)
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
		Name: truncateName(aud.Name + "_" + effect.Name()),
	}

	if !save {
		if err := s.store.Replace(ctx, audioId, newAud, out); err != nil {
			return audio.Audio{}, fmt.Errorf("replace audio: %w", err)
		}
	} else {
		newAud.ID = uuid.NewString()
		if err := s.store.Save(ctx, newAud, out); err != nil {
			return audio.Audio{}, fmt.Errorf("save new audio: %w", err)
		}
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
		effect, err := s.registry.Bind(payload)
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
		Name: truncateName(aud.Name + "_" + applied),
	}

	if !save {
		if err := s.store.Replace(ctx, audioId, newAud, data); err != nil {
			return audio.Audio{}, "", fmt.Errorf("replace audio: %w", err)
		}
	} else {
		newAud.ID = uuid.NewString()
		if err := s.store.Save(ctx, newAud, data); err != nil {
			return audio.Audio{}, "", fmt.Errorf("save new audio: %w", err)
		}
	}

	return newAud, applied, nil

}

func (s *EffectService) ListAvailable() []string {
	return s.registry.Names()
}

// maxNameLen 音频名称长度上限，超长时保留头尾。
const maxNameLen = 64

func truncateName(s string) string {
	if len(s) <= maxNameLen {
		return s
	}
	keep := maxNameLen/2 - 2
	return s[:keep] + "..." + s[len(s)-keep:]
}

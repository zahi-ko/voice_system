package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"voice_system/internal/application/commands"
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
func (s *EffectService) ApplyChain(ctx context.Context, cmd commands.ApplyEffectChain) (audio.Audio, error) {
	// TODO: 复用 Apply 的单步流程，循环展开 Steps；返回链末端音频。
	return audio.Audio{}, errors.New("effect service: ApplyChain not implemented")
}

func (s *EffectService) ListAvailable() []string {
	return s.registry.Names()
}

func (s *EffectService) bindEffect(payload effects.Payload) (effects.Effect, error) {
	effectMap := s.registry.GetMap()
	if effectMap == nil {
		return nil, errors.New("effect registry is nil")
	}

	name := payload["name"].(string)
	eff, ok := effectMap[name]

	if !ok {
		return nil, fmt.Errorf("effect %s not found", name)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	if err := json.Unmarshal(data, &eff); err != nil {
		return nil, fmt.Errorf("unmarshal payload to effect: %w", err)
	}

	return eff, nil

}

package application

import (
	"context"
	"errors"
	"fmt"

	"voice_system/internal/application/commands"
	"voice_system/internal/application/ports"
	"voice_system/internal/domain/audio"
	"voice_system/internal/domain/effects"
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

func (s *EffectService) Apply(ctx context.Context, audioId string, save bool, effect effects.Effect) (audio.Audio, error) {
	// TODO:
	//  1. store.Get 拿到文件流（ReadSeeker），decoder.DecodeMeta 取 sampleRate + Decode 取样本
	//  2. registry.Build(cmd.Step.Name, cmd.Step.Params) 构造 Effect
	//  3. eff.Apply(samples, meta) 得到新样本和元数据
	//  4. encoder.Encode -> saveAsNew ? 新 ID 存储 : 覆盖原文件（覆盖前保留旧版本供 undo）
	//  5. 记录历史（EffectStep + 版本链）

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

package application

import (
	"context"
	"errors"

	"voice_system/internal/application/commands"
	"voice_system/internal/application/ports"
	"voice_system/internal/domain/audio"
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
		last:     "",
	}
}

// Apply 对 audioId 对应的音频应用单个效果。
func (s *EffectService) Apply(ctx context.Context, cmd commands.ApplyEffect) (audio.Audio, error) {
	// TODO:
	//  1. store.Get 拿到文件流（ReadSeeker），decoder.DecodeMeta 取 sampleRate + Decode 取样本
	//  2. registry.Build(cmd.Step.Name, cmd.Step.Params) 构造 Effect
	//  3. eff.Apply(samples, meta) 得到新样本和元数据
	//  4. encoder.Encode -> saveAsNew ? 新 ID 存储 : 覆盖原文件（覆盖前保留旧版本供 undo）
	//  5. 记录历史（EffectStep + 版本链）
	return audio.Audio{}, errors.New("effect service: Apply not implemented")
}

// ApplyChain 按顺序应用效果链，每一步以前一步输出为输入。
func (s *EffectService) ApplyChain(ctx context.Context, cmd commands.ApplyEffectChain) (audio.Audio, error) {
	// TODO: 复用 Apply 的单步流程，循环展开 Steps；返回链末端音频。
	return audio.Audio{}, errors.New("effect service: ApplyChain not implemented")
}

// Undo 撤销最近一次应用，恢复到上一版本音频。
func (s *EffectService) Undo(ctx context.Context, cmd commands.UndoEffect) (audio.Audio, error) {
	// TODO: 依赖版本/历史记录；若采用 saveAsNew 版本链，undo 即回退指针。
	return audio.Audio{}, errors.New("effect service: Undo not implemented")
}

// ListAvailable 返回已注册的效果名列表（ListEffects）。
func (s *EffectService) ListAvailable() []string {
	return s.registry.Names()
}

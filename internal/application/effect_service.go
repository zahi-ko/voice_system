package application

import (
	"context"
	"errors"

	"voice_system/internal/application/commands"
	"voice_system/internal/application/ports"
	"voice_system/internal/domain/audio"
	"voice_system/internal/domain/effects"
)

// EffectService 负责把效果「作用到一份存储的音频」：
// decode -> registry 构造 Effect -> Apply（纯内存变换）-> encode -> store。
// 本层不写任何 DSP 数学；undo 依赖历史版本而非逆运算。
type EffectService struct {
	decoder  ports.AudioDecoder
	encoder  ports.AudioEncoder
	store    ports.AudioStore
	registry *effects.Registry
}

func NewEffectService(
	decoder ports.AudioDecoder,
	encoder ports.AudioEncoder,
	store ports.AudioStore,
	registry *effects.Registry,
) *EffectService {
	return &EffectService{
		decoder:  decoder,
		encoder:  encoder,
		store:    store,
		registry: registry,
	}
}

// Apply 对 audioId 对应的音频应用单个效果。
func (s *EffectService) Apply(ctx context.Context, cmd commands.ApplyEffect) (audio.Audio, error) {
	// TODO:
	//  1. store.Get 拿到文件流（ReadSeeker），decoder.DecodeMeta 取 sampleRate + Decode 取样本
	//  2. registry.Build(cmd.Step.Name, cmd.Step.Params) 构造 Effect
	//  3. eff.Apply(samples, sampleRate) 得到新样本
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

package effects

import (
	"fmt"

	"voice_system/internal/domain/audio"
)

// ParameterNormalization 与 API normalizationParams 对齐（targetLevel: [-60, 0]）。
type ParameterNormalization struct {
	TargetLevel float32 `json:"targetLevel"`
}

type EffectNormalization struct {
	param ParameterNormalization
}

func NewNormalization(p ParameterNormalization) (EffectNormalization, error) {
	if p.TargetLevel < -60 || p.TargetLevel > 0 {
		return EffectNormalization{}, fmt.Errorf("effects: targetLevel %.1fdB out of range [-60, 0]", p.TargetLevel)
	}
	return EffectNormalization{param: p}, nil
}

func (e EffectNormalization) Name() string { return "normalization" }

func (e EffectNormalization) Apply(data audio.AudioData, _ uint32) (audio.AudioData, error) {
	// TODO: 求峰值，按 targetLevel 线性缩放至目标电平，返回新缓冲区。
	return nil, fmt.Errorf("effects: normalization not implemented")
}

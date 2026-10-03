package effects

import (
	"fmt"
	"math"

	"voice_system/internal/domain/audio"
)

type ParameterNormalization struct {
	TargetLevel float32 `json:"targetLevel"`
}

type EffectNormalization struct {
	Param ParameterNormalization `json:"parameters"`
}

func NewNormalization(p ParameterNormalization) (EffectNormalization, error) {
	if p.TargetLevel < -60 || p.TargetLevel > 0 {
		return EffectNormalization{}, fmt.Errorf("effects: targetLevel %.1fdB out of range [-60, 0]", p.TargetLevel)
	}
	return EffectNormalization{Param: p}, nil
}

// NewNormalizationFromPayload 从 JSON 载荷构造归一化效果，走构造函数完成参数校验。
func NewNormalizationFromPayload(p Payload) (Effect, error) {
	var e EffectNormalization
	if err := decodePayload(p, &e); err != nil {
		return nil, err
	}
	return NewNormalization(e.Param)
}

func (e EffectNormalization) Apply(data audio.AudioData, meta audio.Metadata) (audio.AudioData, audio.Metadata, error) {
	out := make(audio.AudioData, len(data))

	peak := math.Inf(-1)
	for _, s := range data {
		val := math.Abs(s)
		if val > peak {
			peak = val
		}
	}

	if peak < 1e-9 {
		copy(out, data)
		return out, meta, nil
	}

	target := math.Pow(10, float64(e.Param.TargetLevel/20.0))
	gain := target / peak
	for i, s := range data {
		out[i] = s * gain
	}

	return out, meta, nil
}

func (e EffectNormalization) Name() string {
	return "normalize"
}

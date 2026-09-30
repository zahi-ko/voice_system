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
	name  string                 `json:"name"`
	Param ParameterNormalization `json:"parameters"`
}

func NewNormalization(p ParameterNormalization) (EffectNormalization, error) {
	if p.TargetLevel < -60 || p.TargetLevel > 0 {
		return EffectNormalization{}, fmt.Errorf("effects: targetLevel %.1fdB out of range [-60, 0]", p.TargetLevel)
	}
	return EffectNormalization{name: "normalize", Param: p}, nil
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
	return e.name
}
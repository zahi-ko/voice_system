package effects

import (
	"fmt"
	"math"

	"voice_system/internal/domain/audio"
)

type ParameterGain struct {
	GainDB float32 `json:"db"`
}

// EffectGain 增益效果。作为 Effect 的参考实现，
// 其余效果按「参数结构体 + 构造校验 + Apply」同构展开。
type EffectGain struct {
	Param ParameterGain `json:"parameters"`
}

func NewGain(p ParameterGain) (EffectGain, error) {
	if p.GainDB < -60 || p.GainDB > 12 {
		return EffectGain{}, fmt.Errorf("effects: gain %.1fdB out of range [-60, 12]", p.GainDB)
	}
	return EffectGain{Param: p}, nil
}

func (e EffectGain) Apply(data audio.AudioData, meta audio.Metadata) (audio.AudioData, audio.Metadata, error) {
	factor := math.Pow(10, float64(e.Param.GainDB/20))
	out := make(audio.AudioData, len(data))
	for i, s := range data {
		out[i] = s * factor
	}
	return out, meta, nil
}

func (e EffectGain) Name() string {
	return "gain"
}

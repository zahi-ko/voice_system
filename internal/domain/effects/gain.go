package effects

import (
	"fmt"
	"math"

	"voice_system/internal/domain/audio"
)

// ParameterGain 与 API gainParams 对齐（db: [-60, 12]）。
type ParameterGain struct {
	GainDB float32 `json:"db"`
}

// EffectGain 增益效果。作为 Effect 的参考实现，
// 其余效果按「参数结构体 + 构造校验 + Apply」同构展开。
type EffectGain struct {
	param ParameterGain
}

func NewGain(p ParameterGain) (EffectGain, error) {
	if p.GainDB < -60 || p.GainDB > 12 {
		return EffectGain{}, fmt.Errorf("effects: gain %.1fdB out of range [-60, 12]", p.GainDB)
	}
	return EffectGain{param: p}, nil
}

func (e EffectGain) Name() string { return "gain" }

func (e EffectGain) Apply(data audio.AudioData, _ uint32) (audio.AudioData, error) {
	factor := math.Pow(10, float64(e.param.GainDB)/20)
	out := make(audio.AudioData, len(data))
	for i, s := range data {
		out[i] = float32(float64(s) * factor)
	}
	return out, nil
}

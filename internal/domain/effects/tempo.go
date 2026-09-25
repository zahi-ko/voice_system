package effects

import (
	"fmt"

	"voice_system/internal/domain/audio"
)

// ParameterTempo 与 API tempoParams 对齐（tempo: [0.5, 2.0]）。
type ParameterTempo struct {
	Tempo float32 `json:"tempo"`
}

// EffectTempo 变速不变调（时间伸缩）。API discriminator: "tempo" / "tempo_ola"。
type EffectTempo struct {
	param ParameterTempo
}

func NewTempo(p ParameterTempo) (EffectTempo, error) {
	if p.Tempo < 0.5 || p.Tempo > 2.0 {
		return EffectTempo{}, fmt.Errorf("effects: tempo %.2f out of range [0.5, 2.0]", p.Tempo)
	}
	return EffectTempo{param: p}, nil
}

func (e EffectTempo) Name() string { return "tempo" }

func (e EffectTempo) Apply(data audio.AudioData, sampleRate uint32) (audio.AudioData, error) {
	// TODO: OLA / WSOLA 时间伸缩，输出长度 = len(data) / tempo。
	return nil, fmt.Errorf("effects: tempo not implemented")
}

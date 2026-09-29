package effects

import (
	"fmt"

	"voice_system/internal/domain/audio"
)

type ParameterTempo struct {
	Tempo float32 `json:"tempo"`
}

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

// 直接修改采样率即可达到变速变调的效果
func (e EffectTempo) Apply(data audio.AudioData, meta audio.Metadata) (audio.AudioData, audio.Metadata, error) {
	out := make(audio.AudioData, len(data))
	copy(out, data)

	sr := float32(meta.SampleRate)
	out_sr := sr * e.param.Tempo

	meta.SampleRate = uint32(out_sr)
	meta.Duration = meta.Duration / e.param.Tempo

	return out, meta, nil
}

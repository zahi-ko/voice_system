package effects

import (
	"voice_system/internal/domain/audio"
)

type EffectReverse struct{}

func NewReverse() EffectReverse { return EffectReverse{} }

func (e EffectReverse) Name() string { return "reverse" }

func (e EffectReverse) Apply(data audio.AudioData, meta audio.Metadata) (audio.AudioData, audio.Metadata, error) {
	// TODO: 样本倒序输出；注意 API oneOf 目前未含 reverse，仅为域内预留。
	out := reversed(data)
	return out, meta, nil
}

func reversed(s audio.AudioData) audio.AudioData {
	out := make(audio.AudioData, len(s))
	for i := range s {
		out[i] = s[len(s)-1-i]
	}
	return out
}

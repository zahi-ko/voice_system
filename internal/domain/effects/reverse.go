package effects

import (
	"voice_system/internal/domain/audio"
)

type EffectReverse struct {
	Param ParameterReverse `json:"parameters"`
}

type ParameterReverse struct{}

func NewReverse() EffectReverse { return EffectReverse{} }

// NewReverseFromPayload 从 JSON 载荷构造倒放效果。
func NewReverseFromPayload(p Payload) (Effect, error) {
	var e EffectReverse
	if err := decodePayload(p, &e); err != nil {
		return nil, err
	}
	return NewReverse(), nil
}

func (e EffectReverse) Apply(data audio.AudioData, meta audio.Metadata) (audio.AudioData, audio.Metadata, error) {
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

func (e EffectReverse) Name() string {
	return "reverse"
}

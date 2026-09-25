package effects

import (
	"fmt"

	"voice_system/internal/domain/audio"
)

type EffectReverse struct{}

func NewReverse() EffectReverse { return EffectReverse{} }

func (e EffectReverse) Name() string { return "reverse" }

func (e EffectReverse) Apply(data audio.AudioData, _ uint32) (audio.AudioData, error) {
	// TODO: 样本倒序输出；注意 API oneOf 目前未含 reverse，仅为域内预留。
	return nil, fmt.Errorf("effects: reverse not implemented")
}

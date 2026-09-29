package ports

import "voice_system/internal/domain/effects"

type EffectRegistry interface {
	Register(name string, effect effects.Effect)
	Get(name string) (effects.Effect, bool)
	Names() []string
}

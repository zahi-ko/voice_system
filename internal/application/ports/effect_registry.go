package ports

import "voice_system/internal/domain/effects"

type EffectRegistry interface {
	// Bind 按 payload 中的 name 构造效果实例并完成参数校验。
	Bind(payload effects.Payload) (effects.Effect, error)
	Names() []string
}

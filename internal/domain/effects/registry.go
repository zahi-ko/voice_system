package effects

import (
	"encoding/json"
	"fmt"
)

// builder 从原始参数 JSON 构造效果实例（含参数校验）。
type builder func(raw json.RawMessage) (Effect, error)

// Registry 维护 name -> builder 的映射，
// 对应 API EffectEntityBase 的 discriminator（reverb/gain/tempo/normalization）。
// transport 只传 name + 原始参数，参数解析与校验收敛在各效果的构造函数里。
type Registry struct {
	builders map[string]builder
}

func NewRegistry() *Registry {
	r := &Registry{builders: make(map[string]builder)}
	r.mustRegister("gain", unmarshalInto(func(p ParameterGain) (Effect, error) { return NewGain(p) }))
	// TODO: 注册 normalization / tempo / reverb / reverse
	return r
}

func (r *Registry) mustRegister(name string, b builder) {
	if _, dup := r.builders[name]; dup {
		panic(fmt.Sprintf("effects: duplicate registry name %q", name))
	}
	r.builders[name] = b
}

// Build 按名称构造效果实例。
func (r *Registry) Build(name string, raw json.RawMessage) (Effect, error) {
	b, ok := r.builders[name]
	if !ok {
		return nil, fmt.Errorf("effects: unknown effect %q", name)
	}
	return b(raw)
}

// Names 返回已注册的效果名（ListEffects 用）。
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.builders))
	for n := range r.builders {
		names = append(names, n)
	}
	return names
}

// unmarshalInto 适配「无参校验失败即返回错误」的构造函数签名。
func unmarshalInto[T any](ctor func(T) (Effect, error)) builder {
	return func(raw json.RawMessage) (Effect, error) {
		var p T
		if err := json.Unmarshal(raw, &p); err != nil {
			return nil, fmt.Errorf("effects: parse parameters: %w", err)
		}
		return ctor(p)
	}
}

package commands

import "encoding/json"

// EffectStep 效果链中的一步：效果名 + 原始参数。
// 参数解析与校验由 domain/effects.Registry.Build 完成，
// 命令层不做任何语义解释。
type EffectStep struct {
	Name   string
	Params json.RawMessage
}

// ApplyEffect 应用单个效果（POST /audios/{audioId}/effects）。
type ApplyEffect struct {
	AudioID   string
	SaveAsNew bool
	Step      EffectStep
}

// ApplyEffectChain 按顺序应用效果链（saveAsNew 时新文件为链末端结果）。
type ApplyEffectChain struct {
	AudioID   string
	SaveAsNew bool
	Steps     []EffectStep
}

// UndoEffect 撤销最近一次应用的效果。
type UndoEffect struct {
	AudioID string
}

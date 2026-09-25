package effects

import "voice_system/internal/domain/audio"

// Effect 是一次纯内存的样本变换。
// 实现约束：
//   - 只依赖 domain/audio，不得引入 ports / 存储 / 文件 IO
//   - Apply 必须无副作用：不修改传入的 data，返回新缓冲区
//
// 应用层（effect_service）负责 decode -> Apply -> encode -> store 的编排，
// 本接口不关心音频从哪来、到哪去。
type Effect interface {
	// Name 返回效果的唯一标识（与 API discriminator 一致，如 "gain"）。
	Name() string

	// Apply 对样本执行变换。sampleRate 为采样率（Hz），时间类效果需要。
	Apply(data audio.AudioData, sampleRate uint32) (audio.AudioData, error)
}

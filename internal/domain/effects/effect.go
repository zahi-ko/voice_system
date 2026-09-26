package effects

import "voice_system/internal/domain/audio"

type Effect interface {
	Name() string

	Apply(data audio.AudioData, sampleRate uint32) (audio.AudioData, error)
}

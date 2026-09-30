package effects

import "voice_system/internal/domain/audio"

type Effect interface {
	Apply(data audio.AudioData, meta audio.Metadata) (audio.AudioData, audio.Metadata, error)
	Name() string
}

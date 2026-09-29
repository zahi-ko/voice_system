package effects

import "voice_system/internal/domain/audio"

type Effect interface {
	Name() string

	Apply(data audio.AudioData, meta audio.Metadata) (audio.AudioData, audio.Metadata, error)
}

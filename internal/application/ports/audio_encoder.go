package ports

import (
	"io"

	"voice_system/internal/domain/audio"
)

type AudioEncoder interface {
	Encode(audio.AudioData, audio.Metadata) (io.ReadSeeker, error)
}

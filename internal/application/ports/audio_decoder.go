package ports

import (
	"io"

	"voice_system/internal/domain/audio"
)

type AudioDecoder interface {
	DecodeMeta(io.ReadSeeker) (audio.Metadata, error)
	Decode(io.ReadSeeker) (audio.AudioData, error)
}

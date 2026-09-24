package audioio

import (
	"io"

	"voice_system/internal/domain/audio"
)

const (
	gomp3NumChannels   = 2
	gomp3Precision     = 2
	gomp3BytesPerFrame = gomp3NumChannels * gomp3Precision
)

type MP3Decoder struct{}

func NewMP3Decoder() *MP3Decoder {
	return &MP3Decoder{}
}

func (d *MP3Decoder) DecodeMeta(input io.ReadSeeker) (audio.Metadata, error) {
	ad, err := audiomorph.DecodeFile()
}

// func (d *MP3Decoder) Decode(input io.ReadSeeker) (audio.AudioData, error) {
// }

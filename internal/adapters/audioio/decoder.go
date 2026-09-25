package audioio

import (
	"errors"
	"fmt"
	"io"

	"voice_system/internal/domain/audio"

	"github.com/zahi-ko/audiomorph"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) DecodeMeta(input io.ReadSeeker) (audio.Metadata, error) {
	if input == nil {
		return audio.Metadata{}, errors.New("audio input is nil")
	}

	metadata, err := audiomorph.DecodeMetadata(input)
	if err != nil {
		return audio.Metadata{}, err
	}

	return audio.Metadata{
		Format:     audio.Format(metadata.Format),
		Duration:   float32(metadata.Duration),
		SampleRate: uint32(metadata.SampleRate),
		BitDepth:   uint8(metadata.BitDepth),
	}, nil
}

func (d *Decoder) Decode(input io.ReadSeeker) (audio.AudioData, error) {
	if input == nil {
		return audio.AudioData{}, errors.New("audio input is nil")
	}

	data, err := audiomorph.Decode(input)
	if err != nil {
		return audio.AudioData{}, fmt.Errorf("decode audio: %w", err)
	}

	return data.Data, nil
}

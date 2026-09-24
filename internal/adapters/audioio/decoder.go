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

// func (d *Decoder) DecodeMeta(input io.ReadSeeker) (audio.Metadata, error) {
// 	if input == nil {
// 		return audio.Metadata{}, errors.New("audio input is nil")
// 	}

// 	a, err := audiomorph.Decode(input)
// 	if err != nil {
// 		return audio.Metadata{}, err
// 	}

// 	return audio.Metadata{
// 		Channels:   uint8(a.NumChannels),
// 		SampleRate: uint32(a.SampleRate),
// 		Duration:   float32(a.Duration),
// 		Format:     audio.Format(a.Format),
// 	}, nil
// }

func (d *Decoder) Decode(input io.ReadSeeker) (audio.AudioData, audio.Metadata, error) {
	if input == nil {
		return audio.AudioData{}, audio.Metadata{}, errors.New("audio input is nil")
	}

	data, err := audiomorph.Decode(input)
	if err != nil {
		return audio.AudioData{}, audio.Metadata{}, fmt.Errorf("decode audio: %w", err)
	}

	if data.NumChannels > 1 {
		mono := make([]float32, len(data.Data)/data.NumChannels)
		for frame := range mono {
			for channel := 0; channel < data.NumChannels; channel++ {
				mono[frame] += data.Data[frame*data.NumChannels+channel]
			}
			mono[frame] /= float32(data.NumChannels)
		}
		data.Data = mono
	}

	return data.Data, audio.Metadata{
		Channels:   1,
		SampleRate: uint32(data.SampleRate),
		Duration:   float32(data.Duration),
		Format:     audio.Format(data.Format),
	}, nil
}

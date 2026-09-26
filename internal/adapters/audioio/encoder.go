package audioio

import (
	"fmt"
	"io"

	"voice_system/internal/domain/audio"

	"github.com/zahi-ko/audiomorph"
)

type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (encoder *Encoder) Encode(data audio.AudioData, metadata audio.Metadata) (io.ReadSeeker, error) {
	a := &audiomorph.Audio{
		SampleRate: int(metadata.SampleRate),
		BitDepth:   int(metadata.BitDepth),
		Format:     string(metadata.Format),
		Data:       data,
		Duration:   float64(metadata.Duration),
	}

	output, err := audiomorph.Encode(a)
	if err != nil {
		return nil, fmt.Errorf("Encode file: %w", err)
	}

	return output, nil
}

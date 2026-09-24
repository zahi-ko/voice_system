package audioio

import (
	"fmt"
	"io"

	"voice_system/internal/domain/audio"

	"github.com/go-audio/wav"

	gaudio "github.com/go-audio/audio"
)

type WAVDecoder struct{}

func NewWAVDecoder() *WAVDecoder {
	return &WAVDecoder{}
}

func (d *WAVDecoder) DecodeMeta(input io.ReadSeeker) (audio.Metadata, error) {
	decoder := wav.NewDecoder(input)
	if !decoder.IsValidFile() {
		return audio.Metadata{}, fmt.Errorf("invalid WAV file")
	}

	duration, err := decoder.Duration()
	if err != nil {
		return audio.Metadata{}, fmt.Errorf("read duration: %w", err)
	}

	return audio.Metadata{
		Channels:   uint8(decoder.NumChans),
		SampleRate: decoder.SampleRate,
		Duration:   float32(duration.Seconds()),
	}, nil
}

func (d *WAVDecoder) Decode(input io.ReadSeeker) (audio.AudioData, error) {
	decoder := wav.NewDecoder(input)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid WAV file")
	}

	bitDepth := int(decoder.BitDepth)
	totalBytes := int(decoder.PCMLen())
	bytesPerSample := bitDepth / 8
	totalSamples := totalBytes / bytesPerSample

	samples := make(audio.AudioData, 0, totalSamples)

	buf := &gaudio.IntBuffer{
		Data: make([]int, 4096),
	}

	for {
		n, err := decoder.PCMBuffer(buf)
		if err != nil {
			return nil, fmt.Errorf("decode file: %w", err)
		} else if n == 0 {
			break
		}

		samples = append(samples, buf.Data[:n]...)
	}

	return samples, nil
}

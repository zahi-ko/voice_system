package audio

import "errors"

type Audio struct {
	ID         string
	Name       string
	Format     Format
	SampleRate uint32
	Duration   float32
}

// Mono only
type AudioData []int

type AudioList struct {
	Items []Audio
}

type Metadata struct {
	Format     Format
	Channels   uint8
	SampleRate uint32
	Duration   float32
}

func New(id, name string, metadata Metadata) (Audio, error) {
	if id == "" || name == "" || metadata.Format == "" || metadata.Channels == 0 || metadata.SampleRate == 0 || metadata.Duration <= 0 {
		return Audio{}, errors.New("invalid audio metadata")
	}

	return Audio{
		ID:         id,
		Name:       name,
		Format:     metadata.Format,
		SampleRate: metadata.SampleRate,
		Duration:   metadata.Duration,
	}, nil
}

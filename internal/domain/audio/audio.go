package audio

import "errors"

type Audio struct {
	ID   string
	Name string
	Meta Metadata
}

// Mono only
type AudioData []float32

type AudioList struct {
	Items []Audio
}

type Metadata struct {
	Format     Format
	SampleRate uint32
	BitDepth   uint8
	Duration   float32
}

func New(id, name string, metadata Metadata) (Audio, error) {
	if id == "" || name == "" || metadata.Format == "" || metadata.SampleRate == 0 || metadata.Duration <= 0 {
		return Audio{}, errors.New("invalid audio metadata")
	}

	return Audio{
		ID:   id,
		Name: name,
		Meta: metadata,
	}, nil
}

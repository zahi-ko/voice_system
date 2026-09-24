package application

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"voice_system/internal/application/ports"
	"voice_system/internal/domain/audio"

	"github.com/google/uuid"
)

const maxUploadSize = 100 << 20

type AudioService struct {
	decoder ports.AudioDecoder
	store   ports.AudioStore
}

func NewAudioService(decoder ports.AudioDecoder, store ports.AudioStore) *AudioService {
	return &AudioService{decoder: decoder, store: store}
}

func (s *AudioService) Upload(ctx context.Context, name string, content io.Reader) (audio.Audio, error) {
	data, err := io.ReadAll(io.LimitReader(content, maxUploadSize+1))
	if err != nil {
		return audio.Audio{}, fmt.Errorf("read upload: %w", err)
	}
	if len(data) == 0 || len(data) > maxUploadSize {
		return audio.Audio{}, fmt.Errorf("upload must be between 1 byte and %d MB", maxUploadSize>>20)
	}
	format, err := audio.DetectFormat(data)
	if err != nil {
		return audio.Audio{}, fmt.Errorf("detect audio format: %w", err)
	}

	metadata, err := s.decoder.DecodeMeta(bytes.NewReader(data))
	if err != nil {
		return audio.Audio{}, fmt.Errorf("decode audio: %w", err)
	}
	metadata.Format = format

	item, err := audio.New(uuid.NewString(), name, metadata)
	if err != nil {
		return audio.Audio{}, err
	}
	if err := s.store.Save(ctx, item, bytes.NewReader(data)); err != nil {
		return audio.Audio{}, fmt.Errorf("save audio: %w", err)
	}

	return item, nil
}

func (s *AudioService) List(ctx context.Context) (audio.AudioList, error) {
	return s.store.List(ctx)
}

package application

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"voice_system/internal/application/ports"
	"voice_system/internal/domain/audio"

	"github.com/google/uuid"
)

const maxUploadSize = 100 << 20

type AudioService struct {
	decoder ports.AudioDecoder
	encoder ports.AudioEncoder
	store   ports.AudioStore
}

func NewAudioService(decoder ports.AudioDecoder, encoder ports.AudioEncoder, store ports.AudioStore) *AudioService {
	return &AudioService{decoder: decoder, encoder: encoder, store: store}
}

func (s *AudioService) Upload(ctx context.Context, name string, content io.ReadSeeker) (audio.Audio, error) {
	metadata, err := s.decoder.DecodeMeta(content)
	if format := metadata.Format; format == "ogg" {
		metadata.Format = "wav"
	}
	if err != nil {
		return audio.Audio{}, fmt.Errorf("decode audio metadata: %w", err)
	}

	data, err := s.decoder.Decode(content)
	if err != nil {
		return audio.Audio{}, fmt.Errorf("decode audio data: %w", err)
	}

	fname := strings.TrimSuffix(name, filepath.Ext(name))
	item, err := audio.New(uuid.NewString(), fname, metadata)
	if err != nil {
		return audio.Audio{}, err
	}
	if err := s.store.Save(ctx, item, data); err != nil {
		return audio.Audio{}, fmt.Errorf("save audio: %w", err)
	}

	return item, nil
}

func (s *AudioService) List(ctx context.Context) (audio.AudioList, error) {
	return s.store.List(ctx)
}

func (s *AudioService) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}

func (s *AudioService) Download(ctx context.Context, id string) (audio.Audio, io.ReadSeeker, error) {
	audioMeta, audioData, err := s.store.Get(ctx, id)
	if err != nil {
		return audio.Audio{}, nil, err
	}

	// 编码为具体文件格式
	data, err := s.encoder.Encode(audioData, audioMeta.Meta)
	if err != nil {
		return audio.Audio{}, nil, fmt.Errorf("encode audio: %w", err)
	}

	return audioMeta, data, nil

}

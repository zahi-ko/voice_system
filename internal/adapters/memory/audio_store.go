package memory

import (
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"voice_system/internal/domain/audio"
)

type AudioStore struct {
	root  string
	mu    sync.RWMutex
	items map[string]audio.Audio
}

func NewAudioStore() (*AudioStore, error) {
	root, err := os.MkdirTemp("", "audio")
	if err != nil {
		return nil, err
	}
	return &AudioStore{root: root, items: make(map[string]audio.Audio)}, nil
}

func (s *AudioStore) Save(ctx context.Context, item audio.Audio, audioData audio.AudioData) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	path := filepath.Join(s.root, item.ID)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create audio file: %w", err)
	}
	defer file.Close()

	// 采用序列化保存[]float32数据
	enc := gob.NewEncoder(file)
	if err := enc.Encode(audioData); err != nil {
		_ = os.Remove(path)
		return fmt.Errorf("write audio file: %w", err)
	}

	s.mu.Lock()
	s.items[item.ID] = item
	s.mu.Unlock()
	return nil
}

func (s *AudioStore) Get(ctx context.Context, id string) (audio.Audio, audio.AudioData, error) {
	select {
	case <-ctx.Done():
		return audio.Audio{}, nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	item, ok := s.items[id]
	s.mu.RUnlock()
	if !ok {
		return audio.Audio{}, nil, fmt.Errorf("audio not found")
	}

	path := filepath.Join(s.root, item.ID)
	file, err := os.Open(path)
	if err != nil {
		return audio.Audio{}, nil, fmt.Errorf("open audio file: %w", err)
	}
	defer file.Close()

	var audioData audio.AudioData
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&audioData); err != nil {
		return audio.Audio{}, nil, fmt.Errorf("read audio file: %w", err)
	}

	return item, audioData, nil
}

func (s *AudioStore) GetInfo(ctx context.Context, id string) (audio.Audio, error) {
	select {
	case <-ctx.Done():
		return audio.Audio{}, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return audio.Audio{}, fmt.Errorf("audio not found")
	}

	return item, nil
}

func (s *AudioStore) List(ctx context.Context) (audio.AudioList, error) {
	select {
	case <-ctx.Done():
		return audio.AudioList{}, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]audio.Audio, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}

	return audio.AudioList{Items: items}, nil
}

func (s *AudioStore) Delete(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, id)
	path := filepath.Join(s.root, id)
	err := os.Remove(path)

	if err != nil {
		return fmt.Errorf("remove file: %w", err)
	}

	return nil
}

package ports

import (
	"context"

	"voice_system/internal/domain/audio"
)

type AudioStore interface {
	Get(context.Context, string) (audio.Audio, audio.AudioData, error)
	GetInfo(context.Context, string) (audio.Audio, error)

	Save(context.Context, audio.Audio, audio.AudioData) error
	Delete(context.Context, string) error

	List(context.Context) (audio.AudioList, error)
}

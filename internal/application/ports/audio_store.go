package ports

import (
	"context"
	"io"

	"voice_system/internal/domain/audio"
)

type AudioStore interface {
	Get(context.Context, string) (audio.Audio, io.ReadCloser, error)
	GetPath(context.Context, string) (string, error)

	Save(context.Context, audio.Audio, io.Reader) error
	Delete(context.Context, string) error

	List(context.Context) (audio.AudioList, error)
}

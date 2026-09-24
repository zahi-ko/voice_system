package ports

import (
	"context"
	"io"

	"voice_system/internal/domain/audio"
)

type AudioStore interface {
	Save(context.Context, audio.Audio, io.Reader) error
	Get(context.Context, string) (audio.Audio, io.ReadCloser, error)
	List(context.Context) (audio.AudioList, error)
}

package http

import (
	"voice_system/internal/domain/audio"
	"voice_system/internal/transport/http/generated"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

func AudioToAPI(item audio.Audio) generated.AudioMeta {
	id := item.ID
	name := item.Name
	sampleRate := generated.AudioMetaSampleRate(item.Meta.SampleRate)
	duration := item.Meta.Duration
	bitDepth := item.Meta.BitDepth

	return generated.AudioMeta{
		ID:         &id,
		Name:       &name,
		SampleRate: &sampleRate,
		Duration:   &duration,
		BitDepth:   int(bitDepth),
	}
}

func AudioListToAPI(items audio.AudioList) generated.AudioList {
	apiItems := make([]generated.AudioMeta, len(items.Items))
	for i, item := range items.Items {
		apiItems[i] = AudioToAPI(item)
	}

	return generated.AudioList{
		Items: apiItems,
	}
}

func APItoUUID(audioId openapi_types.UUID) string {
	return audioId.String()
}

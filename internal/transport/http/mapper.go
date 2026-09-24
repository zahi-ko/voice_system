package http

import (
	"voice_system/internal/domain/audio"
	"voice_system/internal/transport/http/generated"
)

func AudioToAPI(item audio.Audio) generated.AudioMeta {
	id := item.ID
	name := item.Name
	channels := generated.AudioMetaChannels(item.Channels)
	sampleRate := generated.AudioMetaSampleRate(item.SampleRate)
	duration := item.Duration

	return generated.AudioMeta{
		ID:         &id,
		Name:       &name,
		Channels:   &channels,
		SampleRate: &sampleRate,
		Duration:   &duration,
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

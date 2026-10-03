package http

import (
	"uuid"
	"voice_system/internal/domain/audio"
	"voice_system/internal/transport/http/generated"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type EffectChainResponse struct {
	ID             string              `json:"id"`
	Metadata       generated.AudioMeta `json:"metadata"`
	AppliedEffects string              `json:"applied_effects"`
}

func FormatToMIME(format audio.Format) string {
	switch format {
	case "mp3":
		return "audio/mpeg"
	case "wav":
		return "audio/wav"
	case "ogg":
		return "audio/ogg"
	case "flac":
		return "audio/flac"
	case "aiff":
		return "audio/aiff"
	default:
		return "application/octet-stream"
	}
}

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

func UUIDToAPI(id string) openapi_types.UUID {
	converted, err := uuid.Parse(id)
	if err != nil {
		return openapi_types.UUID{}
	}
	return openapi_types.UUID(converted)
}

func APItoSaveAsNew(saveAsNew generated.ApplyEffectParams) bool {
	save := saveAsNew.SaveAsNew
	if save == nil {
		return false
	}
	return *save
}

func APItoSaveAsNewChain(saveAsNew generated.ApplyEffectChainParams) bool {
	save := saveAsNew.SaveAsNew
	if save == nil {
		return false
	}
	return *save
}

func EffectChainResponseToAPI(aud audio.Audio, applied string) EffectChainResponse {
	return EffectChainResponse{
		ID:             aud.ID,
		AppliedEffects: applied,
		Metadata:       AudioToAPI(aud),
	}
}

func APItoPoint(point generated.GetWaveformParams) int {
	if point.Points == nil {
		return 2000
	}
	return *point.Points
}

func APItoNFFT(nfft generated.GetSpectrumParams) int {
	if nfft.Nfft == nil {
		return 1024
	}
	return int(*nfft.Nfft)
}

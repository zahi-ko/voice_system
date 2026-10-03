package effects

import (
	"encoding/json"
	"fmt"

	"voice_system/internal/domain/audio"
)

type Effect interface {
	Apply(data audio.AudioData, meta audio.Metadata) (audio.AudioData, audio.Metadata, error)
	Name() string
}

type Payload map[string]interface{}

// decodePayload 将 JSON 载荷解码到目标结构体。
func decodePayload(payload Payload, out any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("effects: marshal payload: %w", err)
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("effects: unmarshal payload: %w", err)
	}
	return nil
}

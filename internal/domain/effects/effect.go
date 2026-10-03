package effects

import (
	"encoding/json"
	"errors"
	"fmt"

	"voice_system/internal/domain/audio"
)

// ErrInvalidParameters 参数缺失或越界，属于调用方错误，HTTP 层应映射为 400。
var ErrInvalidParameters = errors.New("effects: invalid parameters")

// ErrUnknownEffect 效果名未在注册表中登记，HTTP 层同样映射为 400。
var ErrUnknownEffect = errors.New("effects: unknown effect")

// invalidParameters 用哨兵错误包裹参数校验失败，保留原始文案供响应体使用。
func invalidParameters(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalidParameters, fmt.Sprintf(format, args...))
}

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

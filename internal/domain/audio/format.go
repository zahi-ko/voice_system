package audio

import (
	"bytes"
	"errors"
)

type Format string

const (
	FormatWAV  Format = "wav"
	FormatAIFF Format = "aiff"
	FormatMP3  Format = "mp3"
	FormatOGG  Format = "ogg"
	FormatFLAC Format = "flac"
)

var ErrUnsupportedFormat = errors.New("unsupported audio format")

func DetectFormat(data []byte) (Format, error) {
	switch {
	case len(data) >= 12 && bytes.Equal(data[:4], []byte("RIFF")) && bytes.Equal(data[8:12], []byte("WAVE")):
		return FormatWAV, nil
	case len(data) >= 12 && bytes.Equal(data[:4], []byte("FORM")) &&
		(bytes.Equal(data[8:12], []byte("AIFF")) || bytes.Equal(data[8:12], []byte("AIFC"))):
		return FormatAIFF, nil
	case len(data) >= 4 && bytes.Equal(data[:4], []byte("fLaC")):
		return FormatFLAC, nil
	case len(data) >= 4 && bytes.Equal(data[:4], []byte("OggS")):
		return FormatOGG, nil
	case len(data) >= 3 && bytes.Equal(data[:3], []byte("ID3")):
		return FormatMP3, nil
	case isMP3Frame(data):
		return FormatMP3, nil
	default:
		return "", ErrUnsupportedFormat
	}
}

func isMP3Frame(data []byte) bool {
	if len(data) < 2 || data[0] != 0xff || data[1]&0xe0 != 0xe0 {
		return false
	}

	return data[1]&0x06 != 0
}

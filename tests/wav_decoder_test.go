package audioio

import (
	"os"
	"strings"
	"testing"

	"voice_system/internal/adapters/audioio"

	gaudio "github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func TestWAVDecoderDecodeMeta(t *testing.T) {
	file := createWAVFile(t, 8000, 1, []int{100, -200, 300, -400})

	metadata, err := audioio.NewWAVDecoder().DecodeMeta(file)
	if err != nil {
		t.Fatalf("DecodeMeta() error = %v", err)
	}

	if metadata.Channels != 1 {
		t.Errorf("DecodeMeta() channels = %d, want 1", metadata.Channels)
	}
	if metadata.SampleRate != 8000 {
		t.Errorf("DecodeMeta() sample rate = %d, want 8000", metadata.SampleRate)
	}
	if metadata.Duration <= 0 {
		t.Errorf("DecodeMeta() duration = %f, want a positive duration", metadata.Duration)
	}
}

func TestWAVDecoderDecode(t *testing.T) {
	want := []int{100, -200, 300, -400}
	file := createWAVFile(t, 8000, 1, want)

	got, err := audioio.NewWAVDecoder().Decode(file)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if len(got) != len(want) {
		t.Errorf("Decode() = %v, want %v", got, want)
		return
	}
	for i, sample := range want {
		if got[i] != sample {
			t.Errorf("Decode() sample %d = %d, want %d", i, got[i], sample)
		}
	}
}

func TestWAVDecoderRejectsInvalidFile(t *testing.T) {
	decoder := audioio.NewWAVDecoder()

	if _, err := decoder.DecodeMeta(strings.NewReader("not a WAV file")); err == nil {
		t.Error("DecodeMeta() error = nil, want invalid WAV file error")
	}
	if _, err := decoder.Decode(strings.NewReader("not a WAV file")); err == nil {
		t.Error("Decode() error = nil, want invalid WAV file error")
	}
}

func createWAVFile(t *testing.T, sampleRate, channels int, samples []int) *os.File {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), "decoder-*.wav")
	if err != nil {
		t.Fatalf("create temporary WAV file: %v", err)
	}
	t.Cleanup(func() {
		file.Close()
	})

	encoder := wav.NewEncoder(file, sampleRate, 16, channels, 1)
	buffer := &gaudio.IntBuffer{
		Format: &gaudio.Format{
			SampleRate:  sampleRate,
			NumChannels: channels,
		},
		Data: samples,
	}
	if err := encoder.Write(buffer); err != nil {
		t.Fatalf("write WAV data: %v", err)
	}
	if err := encoder.Close(); err != nil {
		t.Fatalf("close WAV encoder: %v", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("rewind WAV file: %v", err)
	}

	return file
}

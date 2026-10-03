package effects

import (
	"math"
	"testing"

	"voice_system/internal/domain/audio"
)

var testMeta = audio.Metadata{
	Format:     audio.FormatWAV,
	SampleRate: 8000,
	Duration:   1.0,
}

func sineData(n int) audio.AudioData {
	data := make(audio.AudioData, n)
	for i := range data {
		data[i] = math.Sin(2 * math.Pi * 440 * float64(i) / 8000)
	}
	return data
}

func TestGainApply(t *testing.T) {
	cases := []struct {
		name   string
		gainDB float32
	}{
		{"zero gain", 0},
		{"plus six dB", 6},
		{"minus six dB", -6},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e := EffectGain{Param: ParameterGain{GainDB: c.gainDB}}
			in := sineData(100)
			out, meta, err := e.Apply(in, testMeta)
			if err != nil {
				t.Fatalf("apply: %v", err)
			}
			if len(out) != len(in) {
				t.Fatalf("out length = %d, want %d", len(out), len(in))
			}
			if meta.SampleRate != testMeta.SampleRate {
				t.Errorf("metadata must be unchanged by gain")
			}
			factor := math.Pow(10, float64(float32(c.gainDB)/20))
			for i := range in {
				if math.Abs(out[i]-in[i]*factor) > 1e-9 {
					t.Fatalf("sample %d: got %f, want %f", i, out[i], in[i]*factor)
				}
			}
		})
	}
}

func TestGainName(t *testing.T) {
	if got := (EffectGain{}).Name(); got != "gain" {
		t.Errorf("name = %q, want %q", got, "gain")
	}
}

func TestNormalizeApply(t *testing.T) {
	cases := []struct {
		name    string
		levelDB float32
	}{
		{"full scale", 0},
		{"minus six dB", -6},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e := EffectNormalization{Param: ParameterNormalization{TargetLevel: c.levelDB}}
			in := sineData(1000)
			out, meta, err := e.Apply(in, testMeta)
			if err != nil {
				t.Fatalf("apply: %v", err)
			}
			if len(out) != len(in) {
				t.Fatalf("out length = %d, want %d", len(out), len(in))
			}
			if meta.Duration != testMeta.Duration {
				t.Errorf("metadata must be unchanged by normalize")
			}
			peak := 0.0
			for _, s := range out {
				if a := math.Abs(s); a > peak {
					peak = a
				}
			}
			want := math.Pow(10, float64(c.levelDB)/20)
			if math.Abs(peak-want) > 1e-6 {
				t.Errorf("normalized peak = %f, want %f", peak, want)
			}
		})
	}
}

func TestNormalizeSilence(t *testing.T) {
	e := EffectNormalization{Param: ParameterNormalization{TargetLevel: 0}}
	out, _, err := e.Apply(make(audio.AudioData, 10), testMeta)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	for i, s := range out {
		if s != 0 {
			t.Errorf("silence must stay silent, sample %d = %f", i, s)
		}
	}
}

func TestReverseApply(t *testing.T) {
	e := EffectReverse{}
	in := audio.AudioData{1, 2, 3, 4, 5}
	out, meta, err := e.Apply(in, testMeta)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	want := audio.AudioData{5, 4, 3, 2, 1}
	for i := range want {
		if out[i] != want[i] {
			t.Fatalf("out[%d] = %f, want %f", i, out[i], want[i])
		}
	}
	if meta.SampleRate != testMeta.SampleRate {
		t.Errorf("metadata must be unchanged by reverse")
	}
	// 输入不能被原地修改
	if in[0] != 1 || in[4] != 5 {
		t.Errorf("input mutated: %v", in)
	}
}

func TestTempoApply(t *testing.T) {
	cases := []struct {
		name    string
		tempo   float32
		wantSR  uint32
		wantDur float32
	}{
		{"double speed", 2.0, 16000, 0.5},
		{"half speed", 0.5, 4000, 2.0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e := EffectTempo{Param: ParameterTempo{Tempo: c.tempo}}
			in := sineData(100)
			out, meta, err := e.Apply(in, testMeta)
			if err != nil {
				t.Fatalf("apply: %v", err)
			}
			// tempo 变速不变采样数据，仅改采样率与时长元数据
			if len(out) != len(in) {
				t.Fatalf("out length = %d, want %d (samples unchanged)", len(out), len(in))
			}
			for i := range in {
				if out[i] != in[i] {
					t.Fatalf("sample %d changed: %f -> %f", i, in[i], out[i])
				}
			}
			if meta.SampleRate != c.wantSR {
				t.Errorf("sample rate = %d, want %d", meta.SampleRate, c.wantSR)
			}
			if math.Abs(float64(meta.Duration-c.wantDur)) > 1e-6 {
				t.Errorf("duration = %f, want %f", meta.Duration, c.wantDur)
			}
		})
	}
}

func TestEffectNames(t *testing.T) {
	cases := map[Effect]string{
		EffectNormalization{}: "normalize",
		EffectReverse{}:       "reverse",
		EffectTempo{}:         "tempo",
		EffectTempoPV{}:       "tempo_pv",
	}
	for e, want := range cases {
		if got := e.Name(); got != want {
			t.Errorf("name = %q, want %q", got, want)
		}
	}
}

func TestReverseEmptyData(t *testing.T) {
	e := EffectReverse{}
	out, _, err := e.Apply(audio.AudioData{}, testMeta)
	if err != nil {
		t.Fatalf("apply empty: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("out length = %d, want 0", len(out))
	}
}

func TestFromPayloadValidation(t *testing.T) {
	if _, err := NewGainFromPayload(Payload{"parameters": map[string]any{"db": 999.0}}); err == nil {
		t.Error("gain db=999 should be rejected")
	}
	if _, err := NewGainFromPayload(Payload{"parameters": map[string]any{"db": 6.0}}); err != nil {
		t.Errorf("gain db=6 should pass: %v", err)
	}
	if _, err := NewTempoFromPayload(Payload{"parameters": map[string]any{"tempo": 4.0}}); err == nil {
		t.Error("tempo=4.0 should be rejected")
	}
	if _, err := NewTempoFromPayload(Payload{"parameters": map[string]any{"tempo": 1.5}}); err != nil {
		t.Errorf("tempo=1.5 should pass: %v", err)
	}
	if _, err := NewTempoPVFromPayload(Payload{"parameters": map[string]any{"tempo": 0.1}}); err == nil {
		t.Error("tempo_pv tempo=0.1 should be rejected")
	}
	if _, err := NewNormalizationFromPayload(Payload{"parameters": map[string]any{"targetLevel": 3.0}}); err == nil {
		t.Error("targetLevel=3.0 should be rejected")
	}
	if _, err := NewReverseFromPayload(Payload{}); err != nil {
		t.Errorf("reverse should always pass: %v", err)
	}
}

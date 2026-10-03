package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStats(t *testing.T) {
	router, id := newEffectsRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/analysis/"+id+"/stats", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("stats status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var stats struct {
		ID         string  `json:"id"`
		Duration   float32 `json:"duration"`
		SampleRate int     `json:"sampleRate"`
		RMS        float64 `json:"rms"`
		Peak       float64 `json:"peak"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &stats); err != nil {
		t.Fatalf("decode stats response: %v", err)
	}

	if stats.ID != id {
		t.Errorf("stats id = %q, want %q", stats.ID, id)
	}
	if stats.SampleRate <= 0 {
		t.Errorf("sampleRate = %d, want > 0", stats.SampleRate)
	}
	if stats.RMS <= 0 || stats.RMS > 1 {
		t.Errorf("rms = %f, want in (0, 1]", stats.RMS)
	}
	if stats.Peak < stats.RMS || stats.Peak > 1 {
		t.Errorf("peak = %f, want >= rms (%f) and <= 1", stats.Peak, stats.RMS)
	}
}

func TestGetWaveform(t *testing.T) {
	router, id := newEffectsRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/analysis/"+id+"/waveform?points=128", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("waveform status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var wf struct {
		ID       string    `json:"id"`
		Duration float32   `json:"duration"`
		MinVal   []float64 `json:"min_values"`
		MaxVal   []float64 `json:"max_values"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &wf); err != nil {
		t.Fatalf("decode waveform response: %v", err)
	}

	if wf.ID != id {
		t.Errorf("waveform id = %q, want %q", wf.ID, id)
	}
	if len(wf.MinVal) != 128 || len(wf.MaxVal) != 128 {
		t.Fatalf("waveform points: min=%d max=%d, want 128", len(wf.MinVal), len(wf.MaxVal))
	}
	for i := range wf.MinVal {
		if wf.MinVal[i] > 0 || wf.MaxVal[i] < 0 {
			t.Errorf("frame %d: min=%f max=%f, want min<=0<=max", i, wf.MinVal[i], wf.MaxVal[i])
		}
		if math.Abs(wf.MinVal[i]) > 1+1e-6 || math.Abs(wf.MaxVal[i]) > 1+1e-6 {
			t.Errorf("frame %d exceeds float range: min=%f max=%f", i, wf.MinVal[i], wf.MaxVal[i])
		}
	}
}

func TestGetSpectrum(t *testing.T) {
	router, id := newEffectsRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/analysis/"+id+"/spectrum?nfft=512", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("spectrum status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var spec struct {
		ID        string    `json:"id"`
		Data      []float64 `json:"spectrum_data"`
		Frequency []float64 `json:"frequencies"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &spec); err != nil {
		t.Fatalf("decode spectrum response: %v", err)
	}

	nBins := 512/2 + 1
	if spec.ID != id {
		t.Errorf("spectrum id = %q, want %q", spec.ID, id)
	}
	if len(spec.Data) != nBins || len(spec.Frequency) != nBins {
		t.Fatalf("spectrum bins: data=%d freqs=%d, want %d", len(spec.Data), len(spec.Frequency), nBins)
	}
	maxDB := spec.Data[0]
	for _, db := range spec.Data {
		if db > maxDB {
			maxDB = db
		}
	}
	if math.Abs(maxDB) > 1e-9 {
		t.Errorf("spectrum max = %f dB, want 0 (normalized to peak)", maxDB)
	}
	for i, f := range spec.Frequency {
		if f < 0 {
			t.Errorf("frequency[%d] = %f, want >= 0", i, f)
		}
	}
}

func TestAnalysisMissingAudio(t *testing.T) {
	router, _ := newEffectsRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/analysis/00000000-0000-0000-0000-000000000000/stats", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("missing audio stats status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

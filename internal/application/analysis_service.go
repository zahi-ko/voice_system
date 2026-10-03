package application

import (
	"context"
	"fmt"
	"math"
	"slices"
	"voice_system/internal/application/ports"
	"voice_system/internal/domain/analysis"
	"voice_system/internal/domain/audio"

	"gonum.org/v1/gonum/dsp/fourier"
)

type AnalysisService struct {
	store ports.AudioStore
}

func NewAnalysisService(store ports.AudioStore) *AnalysisService {
	return &AnalysisService{
		store: store,
	}
}

func (s *AnalysisService) GetAudioStats(ctx context.Context, audioId string) (analysis.AudioStats, error) {
	aud, data, err := s.store.Get(ctx, audioId)
	if err != nil {
		return analysis.AudioStats{}, fmt.Errorf("get audio: %w", err)
	}

	rms, err := getRMS(data)
	if err != nil {
		return analysis.AudioStats{}, fmt.Errorf("calculate RMS: %w", err)
	}

	peak, err := getPeak(data)
	if err != nil {
		return analysis.AudioStats{}, fmt.Errorf("calculate peak: %w", err)
	}

	stats := analysis.New(
		aud.ID,
		aud.Meta.Duration,
		int(aud.Meta.SampleRate),
		rms,
		peak,
	)

	return stats, nil
}

func (s *AnalysisService) GetWaveForm(ctx context.Context, audioId string, points int) (analysis.Waveform, error) {
	aud, data, err := s.store.Get(ctx, audioId)
	if err != nil {
		return analysis.Waveform{}, fmt.Errorf("get audio: %w", err)
	}

	span := linspaceEdges(0, len(data), points+1)
	maximum := make([]float64, points)
	minimum := make([]float64, points)

	for i := range points {
		segement := data[span[i]:span[i+1]]
		if len(segement) == 0 {
			maximum[i] = .0
			minimum[i] = .0
		} else {
			maximum[i] = maxAbs(segement)
			minimum[i] = -maxAbs(segement)
		}
	}

	return analysis.Waveform{
		ID:       aud.ID,
		Duration: aud.Meta.Duration,
		MaxVal:   maximum,
		MinVal:   minimum,
	}, nil

}

// 返回功率谱的相对分贝值和对应的频率
func (s *AnalysisService) GetSpectrum(ctx context.Context, audioId string, nfft int) (analysis.Spectrum, error) {
	aud, data, err := s.store.Get(ctx, audioId)
	if err != nil {
		return analysis.Spectrum{}, fmt.Errorf("get audio: %w", err)
	}

	if len(data) < nfft {
		padded := make(audio.AudioData, nfft)
		copy(padded, data)
		data = padded
	}

	hop := max(1, nfft/2)

	var frames []audio.AudioData
	for start := 0; start+nfft <= len(data); start += hop {
		frames = append(frames, data[start:start+nfft])
	}

	window := hanningWindow(nfft)
	weight := 0.0
	for _, w := range window {
		weight += w * w
	}

	fft := fourier.NewFFT(nfft)

	nBins := nfft/2 + 1
	powerAcc := make(audio.AudioData, nBins)
	buf := make(audio.AudioData, nfft)

	for _, frame := range frames {
		for i := range frame {
			buf[i] = frame[i] * window[i]
		}

		spec := fft.Coefficients(nil, buf)

		for i, c := range spec {
			p := real(c)*real(c) + imag(c)*imag(c)
			if i > 0 && i < nBins-1 {
				p *= 2
			}
			powerAcc[i] += p
		}
	}

	n := float64(len(frames))
	freqs := make(audio.AudioData, nBins)
	magDB := make(audio.AudioData, nBins)

	sr := float64(aud.Meta.SampleRate)
	norm := n * weight
	for i := range freqs {
		freqs[i] = float64(i) * sr / float64(nfft)
		magDB[i] = 10. * math.Log10(powerAcc[i]/norm+1e-12)
	}

	maxDB := slices.Max(magDB)
	for i := range magDB {
		magDB[i] -= maxDB
	}

	return analysis.Spectrum{
		ID:        aud.ID,
		Data:      magDB,
		Frequency: freqs,
	}, nil
}

func getPeak(data audio.AudioData) (float64, error) {
	peak := maxAbs(data)
	return peak, nil
}

func maxAbs(data audio.AudioData) float64 {
	max := 0.0
	for _, sample := range data {
		if absSample := math.Abs(sample); absSample > max {
			max = absSample
		}
	}
	return max
}

func getRMS(data audio.AudioData) (float64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty audio data")
	}

	sum := 0.0
	for _, sample := range data {
		sum += sample * sample
	}

	return math.Sqrt(sum / float64(len(data))), nil
}

func linspaceEdges(start, stop, edge int) []int {
	edges := make([]int, edge)
	for i := range edge {
		edges[i] = start + (stop-start)*i/edge
	}
	// edges[edge] = stop
	return edges
}

func hanningWindow(n int) []float64 {
	window := make([]float64, n)
	for i := range n {
		window[i] = 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(n-1)))
	}
	return window
}

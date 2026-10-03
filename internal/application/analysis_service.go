package application

import (
	"context"
	"fmt"
	"math"
	"voice_system/internal/application/ports"
	"voice_system/internal/domain/analysis"
	"voice_system/internal/domain/audio"
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
	for i := 0; i <= edge; i++ {
		edges[i] = start + (stop-start)*i/edge
	}
	edges[edge] = stop
	return edges
}

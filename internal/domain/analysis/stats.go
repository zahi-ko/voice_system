package analysis

type AudioStats struct {
	ID         string  `json:"id"`
	Duration   float32 `json:"duration"`
	SampleRate int     `json:"sampleRate"`
	RMS        float64 `json:"rms"`
	Peak       float64 `json:"peak"`
}

func New(id string, duration float32, sampleRate int, rms float64, peak float64) AudioStats {
	return AudioStats{
		ID:         id,
		Duration:   duration,
		SampleRate: sampleRate,
		RMS:        rms,
		Peak:       peak,
	}
}

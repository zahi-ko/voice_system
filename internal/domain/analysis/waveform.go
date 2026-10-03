package analysis

type Waveform struct {
	ID       string    `json:"id"`
	Duration float32   `json:"duration"`
	MinVal   []float64 `json:"min_values"`
	MaxVal   []float64 `json:"max_values"`
}

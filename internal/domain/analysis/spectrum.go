package analysis

type Spectrum struct {
	ID        string    `json:"id"`
	Data      []float64 `json:"spectrum_data"`
	Frequency []float64 `json:"frequencies"`
}

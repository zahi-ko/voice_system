package effects

import (
	"fmt"
	"math"

	"voice_system/internal/domain/audio"

	"gonum.org/v1/gonum/dsp/fourier"
)

type EffectTempoPV struct {
	param ParameterTempo
}

func NewTempoPV(p ParameterTempo) (EffectTempoPV, error) {
	if p.Tempo < 0.5 || p.Tempo > 2.0 {
		return EffectTempoPV{}, fmt.Errorf("effects: tempo %.2f out of range [0.5, 2.0]", p.Tempo)
	}
	return EffectTempoPV{param: p}, nil
}

func (e EffectTempoPV) Name() string { return "tempo_pv" }

func (e EffectTempoPV) Apply(data audio.AudioData, meta audio.Metadata) (audio.AudioData, audio.Metadata, error) {
	pv, err := NewPhaseVoCoder(2048, 512)
	if err != nil {
		return nil, meta, fmt.Errorf("effects: failed to create PhaseVoCoder: %v", err)
	}

	out, err := pv.Stretch(data, e.param.Tempo)
	if err != nil {
		return nil, meta, fmt.Errorf("effects: failed to stretch audio data: %v", err)
	}

	// a := len(out) / meta.SampleRate

	return out, meta, nil
}

type PhaseVoCoder struct {
	FrameSize int
	HopSize   int
	window    []float64
}

func NewPhaseVoCoder(frameSize, hopSize int) (*PhaseVoCoder, error) {
	if frameSize <= 0 || hopSize <= 0 {
		return nil, fmt.Errorf("effects: frameSize and hopSize must be positive")
	}
	window := make([]float64, frameSize)
	for i := range window {
		window[i] = 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(frameSize-1)))
	}
	return &PhaseVoCoder{
		FrameSize: frameSize,
		HopSize:   hopSize,
		window:    window,
	}, nil
}

func (pv *PhaseVoCoder) Stretch(data audio.AudioData, tempo float32) (audio.AudioData, error) {
	if tempo <= 0 {
		return nil, fmt.Errorf("effects: tempo must be positive")
	} else if len(data) == 0 {
		return nil, fmt.Errorf("effects: input audio data is empty")
	}

	speed := float64(tempo)

	n := pv.FrameSize
	ha := pv.HopSize
	hs := max(1, int(math.Round(float64(ha)/speed)))

	nFrames := 1
	if len(data) > n {
		nFrames = int(math.Ceil(float64(len(data)-n)/float64(ha))) + 1
	}

	outLen := (nFrames-1)*hs + n
	out := make(audio.AudioData, outLen)
	weight := make([]float64, outLen)

	frame := make([]float64, n)
	timeFrame := make([]float64, n)

	nFreq := n/2 + 1
	spec := make([]complex128, nFreq)

	prevPhase := make([]float64, nFreq)
	phaseAcc := make([]float64, nFreq)

	fft := fourier.NewFFT(n)

	expectedAdvance := make([]float64, nFreq)
	for k := range nFreq {
		expectedAdvance[k] = 2 * math.Pi * float64(ha) * float64(k) / float64(n)
	}

	for frameIdx := range nFrames {
		start := frameIdx * ha

		for i := range n {
			idx := start + i
			if idx < len(data) {
				frame[i] = data[idx] * pv.window[i]
			} else {
				frame[i] = 0
			}
		}

		_ = fft.Coefficients(spec, frame)

		if frameIdx == 0 {
			for k := range nFreq {
				phase := math.Atan2(imag(spec[k]), real(spec[k]))
				prevPhase[k] = phase
				phaseAcc[k] = phase
			}
		} else {
			for k := range nFreq {
				re := real(spec[k])
				im := imag(spec[k])

				prevRe := math.Cos(prevPhase[k])
				prevIm := math.Sin(prevPhase[k])

				phaseDelta := math.Atan2(im*prevRe-re*prevIm, re*prevRe+im*prevIm)

				expected := expectedAdvance[k]

				deviation := math.Atan2(math.Sin(phaseDelta-expected), math.Cos(phaseDelta-expected))

				actualAdvance := expected + deviation

				phaseAcc[k] += actualAdvance * float64(hs) / float64(ha)
				prevPhase[k] += phaseDelta
			}
		}

		for k := range nFreq {
			mag := math.Hypot(real(spec[k]), imag(spec[k]))
			spec[k] = complex(mag*math.Cos(phaseAcc[k]), mag*math.Sin(phaseAcc[k]))
		}

		_ = fft.Sequence(timeFrame, spec)

		scale := 1.0 / float64(n)
		pos := frameIdx * hs

		for i := range n {
			if pos+i < outLen {
				w := pv.window[i]

				out[pos+i] += timeFrame[i] * w * scale
				weight[pos+i] += w * w
			}
		}
	}

	for i := range out {
		if weight[i] > 1e-9 {
			out[i] /= weight[i]
		}
	}

	targetLen := int(math.Round(float64(len(data)) / speed))
	if targetLen < outLen && targetLen > 0 {
		out = out[:targetLen]
	}

	return out, nil
}

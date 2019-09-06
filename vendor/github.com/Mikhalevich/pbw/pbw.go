package pbw

import (
	"github.com/cheggaaa/pb/v3"
)

type ProgressBarConfig struct {
	Max int64
}

type ProgressBar struct {
	dataChan chan int64
	max      int64
}

func NewProgressBar(c chan int64, cfg ProgressBarConfig) *ProgressBar {
	return &ProgressBar{
		dataChan: c,
		max:      cfg.Max,
	}
}

func (p *ProgressBar) Start() {
	go func() {
		if p.max == 0 {
			p.max = <-p.dataChan
		}

		bar := pb.New64(p.max)
		bar.Start()

		for chunk := range p.dataChan {
			bar.Add64(chunk)
		}
		bar.Finish()
	}()
}

func Show(c chan int64) {
	p := NewProgressBar(c, ProgressBarConfig{})
	p.Start()
}

func ShowWithMax(c chan int64, max int64) {
	p := NewProgressBar(c, ProgressBarConfig{Max: max})
	p.Start()
}

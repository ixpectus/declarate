package output

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Bar struct {
	bar      *progressbar.ProgressBar
	start    time.Time
	finish   time.Time
	duration time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewBar(timeFinish time.Time) *Bar {
	start := time.Now()
	duration := timeFinish.Sub(start)
	ctx, cancel := context.WithCancel(context.Background())
	return &Bar{
		bar:      progressBarDefault(int64(duration.Seconds())),
		start:    start,
		finish:   timeFinish,
		duration: duration,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (b *Bar) Start() {
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	for true {
		select {
		case <-b.ctx.Done():
			b.bar.Exit()
			return
		case <-tick.C:
			b.SetCurrent()
		}
	}
}

func (b *Bar) SetCurrent() {
	b.bar.Set(int(time.Now().Sub(b.start).Seconds()))
}

func (b *Bar) Stop() {
	b.cancel()
	fmt.Fprint(os.Stderr, "\n")
}

func progressBarDefault(max int64) *progressbar.ProgressBar {
	return progressbar.NewOptions64(
		max,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionSpinnerType(2),
		progressbar.OptionShowElapsedTimeOnFinish(),
		// progressbar.OptionSetTheme(progressbar.Theme{}),
		// progressbar.OptionSetPredictTime(false),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowDescriptionAtLineEnd(),
		progressbar.OptionThrottle(65*time.Millisecond),
		// progressbar.OptionShowCount(),
		// progressbar.OptionShowIts(),
		// progressbar.OptionOnCompletion(func() {
		// 	fmt.Fprint(os.Stderr, "\n")
		// }),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)
}

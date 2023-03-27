package run

import "time"

type Poll struct {
	Duration           time.Duration `json:"duration,omitempty" yaml:"duration"`
	Interval           time.Duration `json:"interval,omitempty" yaml:"interval"`
	ResponseBodyRegexp string        `json:"response_body_regexp,omitempty" yaml:"response_body_regexp"`
}

func (p *Poll) PollInterval() []time.Duration {
	if p != nil && p.Duration > 0 {
		interval := p.Interval
		if interval == 0 {
			interval = 1 * time.Second
		}
		duration := p.Duration
		res := []time.Duration{}
		for duration >= 0 && duration >= interval {
			res = append(res, interval)
			duration = duration - interval
		}
		return res
	}
	return []time.Duration{}
}

type Result struct {
	Err      error
	Name     string
	Lvl      int
	FileName string
}

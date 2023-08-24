package run

import (
	"regexp"
	"time"

	"github.com/ixpectus/declarate/compare"
	"github.com/ixpectus/declarate/contract"
)

type Poll struct {
	Duration         time.Duration         `json:"duration,omitempty" yaml:"duration"`
	Interval         time.Duration         `json:"interval,omitempty" yaml:"interval"`
	ResponseRegexp   string                `json:"response_regexp,omitempty" yaml:"response_regexp"`
	ComparisonParams compare.CompareParams `json:"comparisonParams" yaml:"comparisonParams"`
	ResponseTmpls    *string               `json:"response" yaml:"response"`
	comparer         contract.Comparer
}

func (p *Poll) PollInterval() []time.Duration {
	if p != nil && p.Duration > 0 {
		interval := p.Interval
		if interval == 0 {
			interval = 1 * time.Second
		}
		duration := p.Duration
		res := []time.Duration{}
		for duration > 0 && duration >= interval {
			res = append(res, interval)
			duration = duration - interval
		}
		if duration > 0 {
			res = append(res, duration)
		}
		return res
	}
	return []time.Duration{}
}

func (p *Poll) pollContinue(response *string) bool {
	if p.ResponseTmpls != nil && response == nil {
		return false
	}
	if p.ResponseRegexp == "" && p.ResponseTmpls != nil {
		return true
	}
	if p.ResponseRegexp != "" {
		rx, err := regexp.Compile(p.ResponseRegexp)
		if err != nil {
			return false
		}
		if response == nil || !rx.MatchString(*response) {
			return false
		}
		return true
	}
	if p.ResponseTmpls != nil && p.comparer != nil {
		errs, err := p.comparer.CompareJsonBody(*p.ResponseTmpls, *response, p.ComparisonParams)
		if len(errs) > 0 || err != nil {
			return false
		}
	}
	return true
}

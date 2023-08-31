package run

import (
	"regexp"
	"time"

	"github.com/ixpectus/declarate/contract"
)

type Poll struct {
	Duration         time.Duration          `json:"duration,omitempty" yaml:"duration"`
	Interval         time.Duration          `json:"interval,omitempty" yaml:"interval"`
	ResponseRegexp   string                 `json:"response_regexp,omitempty" yaml:"response_regexp"`
	ComparisonParams contract.CompareParams `json:"comparisonParams" yaml:"comparisonParams"`
	ResponseTmpls    *string                `json:"response" yaml:"response"`
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

func (p *Poll) pollContinue(response *string) (bool, []error, error) {
	if p.ResponseRegexp == "" && p.ResponseTmpls == nil {
		return false, nil, nil
	}
	if (p.ResponseTmpls != nil || p.ResponseRegexp != "") && response == nil {
		return false, nil, nil
	}

	if p.ResponseRegexp != "" {
		rx, err := regexp.Compile(p.ResponseRegexp)
		if err != nil {
			return false, nil, nil
		}
		if response == nil || !rx.MatchString(*response) {
			return false, nil, nil
		}
		return true, nil, nil
	}
	if p.ResponseTmpls != nil && p.comparer != nil {
		errs, err := p.comparer.CompareJsonBody(*p.ResponseTmpls, *response, p.ComparisonParams)
		if len(errs) > 0 || err != nil {
			return false, errs, err
		}
	}
	return true, nil, nil
}

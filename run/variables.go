package run

import (
	"github.com/dailymotion/allure-go"
	"github.com/ixpectus/declarate/tools"
	"github.com/ixpectus/declarate/variables"
)

func (r *Runner) fillVariablesByResponse(
	commandResponseBody *string,
	variablesToSet map[string]string,
) error {
	if commandResponseBody == nil || variablesToSet == nil {
		return nil
	}
	jsonVars := map[string]string{}
	for k, v := range variablesToSet {
		if v == "*" {
			r.currentVars.Set(k, *commandResponseBody)
		} else {
			jsonVars[k] = v
		}
	}
	if len(jsonVars) > 0 {
		vars, err := variables.FromJSON(jsonVars, *commandResponseBody, r.currentVars)
		if err != nil {
			return err
		}
		res, _ := r.currentVars.SetAll(vars)
		r.config.Report.AddAttachment("variables from response", allure.TextPlain, []byte(tools.FormatVariables(res)))
	}

	return nil
}

func (r *Runner) fillPersistentVariablesByResponse(
	commandResponseBody *string,
	variablesToSet map[string]string,
) error {
	if commandResponseBody == nil || variablesToSet == nil {
		return nil
	}
	jsonVars := map[string]string{}
	for k, v := range variablesToSet {
		if v == "*" {
		} else {
			jsonVars[k] = v
		}
	}
	if len(jsonVars) > 0 {
		vars, err := variables.FromJSON(jsonVars, *commandResponseBody, r.currentVars)
		if err != nil {
			return err
		}
		res, _ := r.currentVars.SetAll(vars)
		r.config.Report.AddAttachment("variables from response", allure.TextPlain, []byte(tools.FormatVariables(res)))
	}

	return nil
}

func (r *Runner) fillAllVariables(
	commandResponseBody *string,
	conf runConfig,
) error {
	if err := r.fillVariablesByResponse(
		commandResponseBody,
		conf.Variables,
	); err != nil {
		return err
	}
	if err := r.fillPersistentVariablesByResponse(
		commandResponseBody,
		conf.VariablesPersistent,
	); err != nil {
		return err
	}
	return nil
}

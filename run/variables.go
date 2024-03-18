package run

import "github.com/ixpectus/declarate/variables"

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
		for k, v := range vars {
			r.currentVars.Set(k, v)
		}
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
		for k, v := range vars {
			r.currentVars.SetPersistent(k, v)
		}
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

package condition

import "github.com/ixpectus/declarate/contract"

func IsTrue(variables contract.Vars, condition string) bool {
	condition = "$(" + condition + ")"
	condition = variables.Apply(condition)
	return condition == "true"
}

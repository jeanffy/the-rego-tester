package usecases

import (
	"encoding/json"
	"fmt"
	"regotest/internal/domain/ports"
	"regotest/internal/infra/colors"
)

// ---------------------------------------------------------------------------
// #region definition

type RunTestUseCase struct {
	opaPort ports.OpaPort
}

// #endregion

// ---------------------------------------------------------------------------
// #region constructor

type RunTestUseCaseDependencies struct {
	OpaPort ports.OpaPort
}

func NewRunTestUseCase(dependencies RunTestUseCaseDependencies) *RunTestUseCase {
	return &RunTestUseCase{
		opaPort: dependencies.OpaPort,
	}
}

// #endregion

// ---------------------------------------------------------------------------
// #region public

type RunTestUseCaseExecuteParams struct {
	Name       string
	RegoPath   string
	DataPath   string
	EntryPoint string
	Input      json.RawMessage
	Expected   interface{}
	Verbose    bool
}

func (x *RunTestUseCase) Execute(params RunTestUseCaseExecuteParams) bool {
	if params.Verbose {
		fmt.Printf(colors.Dimmed("> Running test '%s'\n"), params.Name)
	}

	result, err := x.opaPort.Evaluate(ports.OpaPortEvaluateParams{
		RegoPath:   params.RegoPath,
		DataPath:   params.DataPath,
		EntryPoint: params.EntryPoint,
		Input:      params.Input,
		Verbose:    params.Verbose,
	})

	if params.Verbose {
		fmt.Printf(colors.Dimmed("> Test output:\n"))
		fmt.Printf(colors.Dimmed("%s\n"), result.Output)
	}

	if err != nil {
		fmt.Printf("❌ [KO] %s\n", params.Name)
		fmt.Printf("evaluate error: %v\n", err)
		return false
	}

	if result.Value != params.Expected {
		fmt.Printf("❌ [KO] %s\n", params.Name)
		fmt.Printf("expected: %v\n", params.Expected)
		fmt.Printf("actual: %v\n", result.Value)
		return false
	}

	fmt.Printf("✅ [OK] %s\n", params.Name)

	return true
}

// #endregion

// ---------------------------------------------------------------------------
// #region events

// #endregion

// ---------------------------------------------------------------------------
// #region private

// #endregion

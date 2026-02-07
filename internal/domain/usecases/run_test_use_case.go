package usecases

import (
	"encoding/json"
	"fmt"
	"theregotester/internal/domain/ports"
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
	result, err := x.opaPort.Evaluate(ports.OpaPortEvaluateParams{
		RegoPath:   params.RegoPath,
		DataPath:   params.DataPath,
		EntryPoint: params.EntryPoint,
		Input:      params.Input,
	})
	if err != nil {
		fmt.Printf("❌ [KO] %s\n", params.Name)
		fmt.Printf("evaluate error: %v\n", err)
		if params.Verbose {
			fmt.Printf("\x1b[2m%s\x1b[0m\n", result.Output)
		}
		return false
	}

	if result.Value != params.Expected {
		fmt.Printf("❌ [KO] %s\n", params.Name)
		fmt.Printf("expected: %v\n", params.Expected)
		fmt.Printf("actual: %v\n", result.Value)
		if params.Verbose {
			fmt.Printf("\x1b[2m%s\x1b[0m\n", result.Output)
		}
		return false
	}

	fmt.Printf("✅ [OK] %s\n", params.Name)
	if params.Verbose {
		fmt.Printf("\x1b[2m%s\x1b[0m\n", result.Output)
	}

	return true
}

// #endregion

// ---------------------------------------------------------------------------
// #region events

// #endregion

// ---------------------------------------------------------------------------
// #region private

// #endregion

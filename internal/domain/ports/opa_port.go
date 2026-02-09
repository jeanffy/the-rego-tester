package ports

import "encoding/json"

var OPA_PORT_TOKEN = "OpaPort"

type OpaPortEvaluateParams struct {
	RegoPath   string
	DataPath   string
	EntryPoint string
	Input      json.RawMessage
	Verbose    bool
}

type OpaPortEvaluateResult struct {
	Value  interface{}
	Output string
}

type OpaPort interface {
	Evaluate(params OpaPortEvaluateParams) (OpaPortEvaluateResult, error)
}

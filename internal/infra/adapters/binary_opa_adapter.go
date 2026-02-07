package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"theregotester/internal/domain/ports"
)

// ---------------------------------------------------------------------------
// #region definition

var _ ports.OpaPort = (*BinaryOpaAdapter)(nil)

type BinaryOpaAdapter struct {
	opaExePath string
}

// #endregion

// ---------------------------------------------------------------------------
// #region constructor

func NewBinaryOpaAdapter() *BinaryOpaAdapter {
	opaPath := filepath.Join(os.TempDir(), "opa")
	return &BinaryOpaAdapter{
		opaExePath: opaPath,
	}
}

// #endregion

// ---------------------------------------------------------------------------
// #region public

type OpaResponse struct {
	Results []OpaResult `json:"result"`
}

type OpaResult struct {
	Expressions []OpaExpression `json:"expressions"`
}

type OpaExpression struct {
	Value    interface{}     `json:"value"`
	Text     string          `json:"text"`
	Location json.RawMessage `json:"location"`
}

func (x *BinaryOpaAdapter) Evaluate(params ports.OpaPortEvaluateParams) (ports.OpaPortEvaluateResult, error) {
	_, err := os.Stat(x.opaExePath)
	if err != nil {
		var url string
		switch fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH) {
		case "darwin-amd64":
			url = "https://github.com/open-policy-agent/opa/releases/download/v1.13.1/opa_darwin_amd64_static"
		case "darwin-arm64":
			url = "https://github.com/open-policy-agent/opa/releases/download/v1.13.1/opa_darwin_arm64_static"
		case "linux-amd64":
			url = "https://github.com/open-policy-agent/opa/releases/download/v1.13.1/opa_linux_amd64_static"
		case "linux-arm64":
			url = "https://github.com/open-policy-agent/opa/releases/download/v1.13.1/opa_linux_arm64_static"
		default:
			panic(fmt.Errorf("not handled GOOS %s", runtime.GOOS))
		}
		fmt.Printf("💡 downloading opa binary from %s\n", url)
		fmt.Printf("   writing binary in %s\n", x.opaExePath)
		if err := downloadOPA(url, x.opaExePath); err != nil {
			panic(err)
		}
	}

	args := []string{
		"eval",
		"--fail",
		"--stdin-input",
		"-d", params.RegoPath,
		"-d", params.DataPath,
		fmt.Sprintf("data.%s", params.EntryPoint),
	}

	inputBytes, err := json.Marshal(params.Input)
	if err != nil {
		return ports.OpaPortEvaluateResult{}, fmt.Errorf("error marshal input: %w", err)
	}

	cmd := exec.Command(x.opaExePath, args...)
	cmd.Stdin = bytes.NewReader(inputBytes)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return ports.OpaPortEvaluateResult{}, fmt.Errorf("opa failed: %v\n", err)
	}

	responseStr, err := extractLastObject(string(out))
	if err != nil {
		panic(err)
	}

	var response OpaResponse
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return ports.OpaPortEvaluateResult{}, fmt.Errorf("invalid json: %w\n", err)
	}

	if len(response.Results) == 0 {
		return ports.OpaPortEvaluateResult{}, fmt.Errorf("no result in response\n")
	}

	result := response.Results[0]

	if len(result.Expressions) == 0 {
		return ports.OpaPortEvaluateResult{}, fmt.Errorf("no expression in result\n")
	}

	expression := result.Expressions[0]

	return ports.OpaPortEvaluateResult{
		Value:  expression.Value,
		Output: string(out),
	}, nil
}

// #endregion

// ---------------------------------------------------------------------------
// #region events

// #endregion

// ---------------------------------------------------------------------------
// #region private

func downloadOPA(url, dest string) error {
	// Create request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	// Create temp file then atomically move into place
	tmp := dest + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	// ensure tmp closed on error
	defer func() {
		out.Close()
		if err != nil {
			os.Remove(tmp)
		}
	}()

	// Copy body to file
	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}
	if err = out.Close(); err != nil {
		return err
	}

	// Make executable
	if err = os.Chmod(tmp, 0o755); err != nil {
		return err
	}

	// Move into final location (overwrite if exists)
	if err = os.Rename(tmp, dest); err != nil {
		// fallback to copy if rename across filesystems fails
		in, err2 := os.Open(tmp)
		if err2 != nil {
			return err2
		}
		defer in.Close()
		fout, err2 := os.Create(dest)
		if err2 != nil {
			return err2
		}
		if _, err2 = io.Copy(fout, in); err2 != nil {
			fout.Close()
			return err2
		}
		fout.Close()
		if err2 = os.Chmod(dest, 0o755); err2 != nil {
			return err2
		}
		os.Remove(tmp)
	}
	return nil
}

// opa reponse is like:
//
// ...
// ... several lines, with print instructions of rego file
// ...
// {
//   "result": [
//     {
//       "expressions": [
//         ...
//       ]
//   ]
// }
//
// last JSON object represents the result of the evaluation

func extractLastObject(s string) (string, error) {
	lines := strings.Split(s, "\n")

	// find last line that starts with '}' (closing) and its index
	endIndex := -1
	for i := len(lines) - 1; i >= 0; i-- {
		if len(lines[i]) > 0 && lines[i][0] == '}' {
			endIndex = i
			break
		}
	}
	if endIndex == -1 {
		return "", fmt.Errorf("no closing '}' found")
	}

	// find the nearest previous line that starts with '{' (opening)
	startIndex := -1
	for i := endIndex; i >= 0; i-- {
		if len(lines[i]) > 0 && lines[i][0] == '{' {
			startIndex = i
			break
		}
	}
	if startIndex == -1 {
		return "", fmt.Errorf("no opening '{' found")
	}

	// join the lines between startIdx and endIdx (inclusive)
	raw := strings.Join(lines[startIndex:endIndex+1], "\n")

	return raw, nil
}

// #endregion

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regotest/internal/domain/ports"
	"regotest/internal/domain/usecases"
	"regotest/internal/infra/adapters"
	"regotest/pkg/di"
)

// #region Test suite

type TestSuite struct {
	Source      Source       `json:"source"`
	EntryPoints []EntryPoint `json:"entryPoints"`
}

type Source struct {
	Rego    string `json:"rego"`
	Data    string `json:"data"`
	Package string `json:"package"`
}

type EntryPoint struct {
	Var   string `json:"var"`
	Tests []Test `json:"tests"`
}

type Test struct {
	Name     string          `json:"name"`
	Input    json.RawMessage `json:"input"`
	Expected interface{}     `json:"expected"`
}

// #endregion Test suite

var VERSION = "0.0.2"

func main() {
	fmt.Printf("RegoTest v%s\n", VERSION)

	verbose := flag.Bool("verbose", false, "verbose output")
	bail := flag.Bool("bail", false, "exit after the first failed test")
	only := flag.String("only", "", "run only one test (name)")
	flag.Parse()
	pos := filter(flag.Args(), func(a string) bool {
		return len(a) > 0
	})

	if len(pos) == 0 {
		fmt.Println("need a test suite file path")
		return
	}

	testSuitePath := pos[0]

	testSuiteContent, err := os.ReadFile(testSuitePath)
	if err != nil {
		fmt.Printf("read %s error: %v\n", testSuitePath, err)
		return
	}

	var testSuite TestSuite
	if err := json.Unmarshal(testSuiteContent, &testSuite); err != nil {
		fmt.Printf("could not marshal json: %v\n", err)
		return
	}

	initDI()

	runTestUseCase := usecases.NewRunTestUseCase(usecases.RunTestUseCaseDependencies{
		OpaPort: di.GetBasicDI().Resolve(ports.OPA_PORT_TOKEN).(ports.OpaPort),
	})

	totalNumberOfTests := 0
	for _, entryPoint := range testSuite.EntryPoints {
		for range entryPoint.Tests {
			totalNumberOfTests++
		}
	}
	numberOfRunTests := 0
	numberOfSucceeded := 0
	numberOfFailed := 0

	for _, entryPoint := range testSuite.EntryPoints {
		for _, test := range entryPoint.Tests {
			if only != nil && *only != "" && test.Name != *only {
				continue
			}

			succeeded := runTestUseCase.Execute(usecases.RunTestUseCaseExecuteParams{
				Name:       test.Name,
				RegoPath:   testSuite.Source.Rego,
				DataPath:   testSuite.Source.Data,
				EntryPoint: fmt.Sprintf("%s.%s", testSuite.Source.Package, entryPoint.Var),
				Input:      test.Input,
				Expected:   test.Expected,
				Verbose:    *verbose,
			})

			numberOfRunTests++
			if succeeded {
				numberOfSucceeded++
			} else {
				numberOfFailed++
				if *bail {
					fmt.Printf("%d tests, %d run, %d succeeded, %d failed\n", totalNumberOfTests, numberOfRunTests, numberOfSucceeded, numberOfFailed)
					os.Exit(1)
				}
			}
		}
	}

	fmt.Printf("%d tests, %d run, %d succeeded, %d failed\n", totalNumberOfTests, numberOfRunTests, numberOfSucceeded, numberOfFailed)

	if numberOfFailed > 0 {
		os.Exit(1)
	}
}

func initDI() {
	di := di.GetBasicDI()
	di.Provide(ports.OPA_PORT_TOKEN, adapters.NewBinaryOpaAdapter())
}

func filter(slice []string, condition func(string) bool) []string {
	var result []string
	for _, value := range slice {
		if condition(value) {
			result = append(result, value)
		}
	}
	return result
}

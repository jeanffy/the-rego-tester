# RegoTest

## Configuration

Test suite JSON file:

```json
{
  "source": {
    "rego": "absolute or relative path to .rego file", // if relative, relative to cwd
    "data": "absolute or relative path to .json data file", // if relative, relative to cwd
    "package": "name of package in .rego file"
  },
  "entryPoints": [ // list of .rego file entry points
    {
      "var": "name of entry point variable",
      "tests": [ // list of test cases to run
        {
          "name": "display name of test",
          "input": { // input object passed to .rego rule (variable 'input' in .rego rule)
            ... // free structure
          },
          // could be any type
          "expected": "expected value to be returned by .rego rule"
        }
      ]
    }
  ]
}
```

## Run

```sh
regotest [-verbose] [-bail] [-only "name"] /path/to/test-suite.json
```

- `-verbose`: enable verbose mode
- `-bail`: exit immediately if a test case fails (don't run subsequent test cases)
- `-only "name"`: run only the test case named "name"

## Build

```sh
sh scripts/build.sh
```

Builds are created in the `dist` directory.

## Debug

Debug configuration is present for VSCode, just add some breakpoint and hit F5.

package generator

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strings"
)

func GenerateTests(node *ast.File, inputFile string) {
	testFileName := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_test.go"
	testFile, err := os.Create(testFileName)
	if err != nil {
		fmt.Printf("Failed to create test file: %v\n", err)
		os.Exit(1)
	}
	defer testFile.Close()

	// Write package declaration
	fmt.Fprintf(testFile, "package %s\n\n", node.Name.Name)
	fmt.Fprintf(testFile, "import (\n")
	fmt.Fprintf(testFile, "\t\"testing\"\n")
	fmt.Fprintf(testFile, "\t\"math\"\n")
	fmt.Fprintf(testFile, "\t\"fmt\"\n")
	fmt.Fprintf(testFile, ")\n\n")

	// Generate test functions
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Recv == nil {
			generateTestFunction(testFile, funcDecl)
		}
	}

	fmt.Printf("Generated test file: %s\n", testFileName)
}

func generateTestFunction(testFile *os.File, funcDecl *ast.FuncDecl) {
	funcName := funcDecl.Name.Name
	testFuncName := "Test" + funcName

	fmt.Fprintf(testFile, "func %s(t *testing.T) {\n", testFuncName)

	// Generate edge cases
	paramNames, paramTypes, edgeCases := getParametersAndEdgeCases(funcDecl.Type.Params)

	// Handle functions with multiple return values
	returnTypes := getReturnTypes(funcDecl.Type.Results)
	multipleReturns := len(returnTypes) > 1

	// Begin test cases
	fmt.Fprintf(testFile, "\ttestCases := []struct {\n")
	for i, paramName := range paramNames {
		fmt.Fprintf(testFile, "\t\t%s %s\n", paramName, paramTypes[i])
	}
	if multipleReturns {
		for i, retType := range returnTypes {
			fmt.Fprintf(testFile, "\t\texpected%d %s\n", i, retType)
		}
	} else {
		fmt.Fprintf(testFile, "\t\texpected %s\n", returnTypes[0])
	}

	// Handle potential errors
	if returnsError(returnTypes) {
		fmt.Fprintf(testFile, "\t\texpectError bool\n")
	}

	fmt.Fprintf(testFile, "\t}{\n")

	// Generate test cases
	for _, params := range edgeCases {
		expected, canCompute := computeExpected(funcName, params)
		fmt.Fprintf(testFile, "\t\t{")
		for _, param := range params {
			fmt.Fprintf(testFile, "%s, ", param)
		}
		if canCompute {
			fmt.Fprintf(testFile, "%s", expected)
			if returnsError(returnTypes) {
				fmt.Fprintf(testFile, ", false")
			}
			fmt.Fprintf(testFile, "},\n")
		} else {
			// Use zero values or placeholders
			if multipleReturns {
				for _, retType := range returnTypes {
					if retType == "error" {
						fmt.Fprintf(testFile, "nil, ")
					} else {
						fmt.Fprintf(testFile, "%s, ", getZeroValue(retType))
					}
				}
				if returnsError(returnTypes) {
					if funcName == "Divide" && containsZero(params[1]) {
						fmt.Fprintf(testFile, "true")
					} else {
						fmt.Fprintf(testFile, "false")
					}
				}
				fmt.Fprintf(testFile, "},\n")
			} else {
				fmt.Fprintf(testFile, "%s", getZeroValue(returnTypes[0]))
				if returnsError(returnTypes) {
					if funcName == "Divide" && containsZero(params[1]) {
						fmt.Fprintf(testFile, ", true")
					} else {
						fmt.Fprintf(testFile, ", false")
					}
				}
				fmt.Fprintf(testFile, "},\n")
			}
		}
	}
	fmt.Fprintf(testFile, "\t}\n\n")

	// Write the test function logic
	fmt.Fprintf(testFile, "\tfor _, tc := range testCases {\n")
	if multipleReturns {
		fmt.Fprintf(testFile, "\t\t")
		for i := range returnTypes {
			if i > 0 {
				fmt.Fprintf(testFile, ", ")
			}
			fmt.Fprintf(testFile, "result%d", i)
		}
		fmt.Fprintf(testFile, " := %s(", funcName)
	} else {
		fmt.Fprintf(testFile, "\t\tresult := %s(", funcName)
	}
	for i, paramName := range paramNames {
		if i > 0 {
			fmt.Fprintf(testFile, ", ")
		}
		fmt.Fprintf(testFile, "tc.%s", paramName)
	}
	fmt.Fprintf(testFile, ")\n")

	// Assertions
	if multipleReturns {
		if returnsError(returnTypes) {
			// Handle error checking
			lastIndex := len(returnTypes) - 1
			fmt.Fprintf(testFile, "\t\tif (result%d != nil) != tc.expectError {\n", lastIndex)
			fmt.Fprintf(testFile, "\t\t\tt.Errorf(\"%s error = %%v, expectError %%v\", result%d, tc.expectError)\n", funcName, lastIndex)
			fmt.Fprintf(testFile, "\t\t}\n")
			// Check other return values if no error
			fmt.Fprintf(testFile, "\t\tif !tc.expectError {\n")
			for i := 0; i < lastIndex; i++ {
				fmt.Fprintf(testFile, "\t\t\tif result%d != tc.expected%d {\n", i, i)
				fmt.Fprintf(testFile, "\t\t\t\tt.Errorf(\"%s expected %%v, got %%v\", tc.expected%d, result%d)\n", funcName, i, i)
				fmt.Fprintf(testFile, "\t\t\t}\n")
			}
			fmt.Fprintf(testFile, "\t\t}\n")
		} else {
			for i := range returnTypes {
				fmt.Fprintf(testFile, "\t\tif result%d != tc.expected%d {\n", i, i)
				fmt.Fprintf(testFile, "\t\t\tt.Errorf(\"%s expected %%v, got %%v\", tc.expected%d, result%d)\n", funcName, i, i)
				fmt.Fprintf(testFile, "\t\t}\n")
			}
		}
	} else {
		if returnsError(returnTypes) {
			fmt.Fprintf(testFile, "\t\tif (err != nil) != tc.expectError {\n")
			fmt.Fprintf(testFile, "\t\t\tt.Errorf(\"%s error = %%v, expectError %%v\", err, tc.expectError)\n", funcName)
			fmt.Fprintf(testFile, "\t\t}\n")
			fmt.Fprintf(testFile, "\t\tif !tc.expectError {\n")
			fmt.Fprintf(testFile, "\t\t\tif result != tc.expected {\n")
			fmt.Fprintf(testFile, "\t\t\t\tt.Errorf(\"%s expected %%v, got %%v\", tc.expected, result)\n", funcName)
			fmt.Fprintf(testFile, "\t\t\t}\n")
			fmt.Fprintf(testFile, "\t\t}\n")
		} else {
			fmt.Fprintf(testFile, "\t\tif result != tc.expected {\n")
			fmt.Fprintf(testFile, "\t\t\tt.Errorf(\"%s expected %%v, got %%v\", tc.expected, result)\n", funcName)
			fmt.Fprintf(testFile, "\t\t}\n")
		}
	}
	fmt.Fprintf(testFile, "\t}\n")
	fmt.Fprintf(testFile, "}\n\n")
}

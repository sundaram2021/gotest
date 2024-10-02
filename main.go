// main.go
package main

import (
    "flag"
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "math"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

func main() {
    // Parse command-line arguments
    var inputFile string
    flag.StringVar(&inputFile, "file", "", "Path to the Go source file")
    flag.Parse()

    if inputFile == "" {
        fmt.Println("Please provide a Go source file using the -file flag.")
        os.Exit(1)
    }

    // Lexical analysis and parsing
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, inputFile, nil, 0)
    if err != nil {
        fmt.Printf("Failed to parse file: %v\n", err)
        os.Exit(1)
    }

    // Generate tests
    generateTests(node, inputFile)
}

func generateTests(node *ast.File, inputFile string) {
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

func getParametersAndEdgeCases(fieldList *ast.FieldList) ([]string, []string, [][]string) {
    var paramNames []string
    var paramTypes []string
    var paramEdgeCases [][]string

    if fieldList == nil {
        return paramNames, paramTypes, [][]string{{}}
    }

    for _, field := range fieldList.List {
        typ := getTypeAsString(field.Type)
        cases := getEdgeCases(typ)
        for _, name := range field.Names {
            paramNames = append(paramNames, name.Name)
            paramTypes = append(paramTypes, typ)
            paramEdgeCases = append(paramEdgeCases, cases)
        }
    }

    edgeCases := cartesianProduct(paramEdgeCases)
    return paramNames, paramTypes, edgeCases
}

func getReturnTypes(fieldList *ast.FieldList) []string {
    var returnTypes []string
    if fieldList == nil {
        return returnTypes
    }
    for _, field := range fieldList.List {
        typ := getTypeAsString(field.Type)
        returnTypes = append(returnTypes, typ)
    }
    return returnTypes
}

func getTypeAsString(expr ast.Expr) string {
    switch t := expr.(type) {
    case *ast.Ident:
        return t.Name
    case *ast.ArrayType:
        return "[]" + getTypeAsString(t.Elt)
    case *ast.StarExpr:
        return "*" + getTypeAsString(t.X)
    case *ast.SelectorExpr:
        return getTypeAsString(t.X) + "." + t.Sel.Name
    case *ast.MapType:
        return "map[" + getTypeAsString(t.Key) + "]" + getTypeAsString(t.Value)
    default:
        return ""
    }
}

func getEdgeCases(typ string) []string {
    switch typ {
    case "int":
        return []string{"0", "1", "-1", "math.MaxInt", "math.MinInt"}
    // ... handle other types similarly
    default:
        if strings.HasPrefix(typ, "[]") {
            elemType := typ[2:]
            elemCases := getEdgeCases(elemType)
            slices := []string{"nil", typ + "{}", typ + "{" + elemCases[0] + "}"}
            return slices
        }
        // For unsupported types, return a zero value or nil
        return []string{getZeroValue(typ)}
    }
}

func getZeroValue(typ string) string {
    switch typ {
    case "int", "int8", "int16", "int32", "int64":
        return "0"
    case "uint", "uint8", "uint16", "uint32", "uint64":
        return "0"
    case "float32", "float64":
        return "0.0"
    case "complex64", "complex128":
        return "complex(0,0)"
    case "string":
        return `""`
    case "bool":
        return "false"
    default:
        if strings.HasPrefix(typ, "[]") {
            return "nil"
        }
        return "nil"
    }
}

func cartesianProduct(slices [][]string) [][]string {
    if len(slices) == 0 {
        return [][]string{{}}
    }

    result := [][]string{{}}

    for _, slice := range slices {
        var newResult [][]string
        for _, res := range result {
            for _, item := range slice {
                newRes := append([]string{}, res...)
                newRes = append(newRes, item)
                newResult = append(newResult, newRes)
            }
        }
        result = newResult
    }

    return result
}

func computeExpected(funcName string, params []string) (string, bool) {
    // For known functions, compute the expected result
    switch funcName {
    case "Add", "Subtract", "Multiply", "Divide":
        // Try to parse the parameters as integers
        aStr := params[0]
        bStr := params[1]
        a, err1 := parseInt(aStr)
        b, err2 := parseInt(bStr)
        if err1 != nil || err2 != nil {
            return "", false
        }
        var result int
        switch funcName {
        case "Add":
            result = a + b
        case "Subtract":
            result = a - b
        case "Multiply":
            result = a * b
        case "Divide":
            if b == 0 {
                return "", false
            }
            result = a / b
        }
        return fmt.Sprintf("%d", result), true
    default:
        // Cannot compute expected result
        return "", false
    }
}

func parseInt(value string) (int, error) {
    value = strings.TrimSpace(value)
    switch value {
    case "math.MaxInt":
        return math.MaxInt, nil
    case "math.MinInt":
        return math.MinInt, nil
    default:
        return strconv.Atoi(value)
    }
}

func returnsError(returnTypes []string) bool {
    for _, typ := range returnTypes {
        if typ == "error" {
            return true
        }
    }
    return false
}

func containsZero(value string) bool {
    trimmed := strings.TrimSpace(value)
    return trimmed == "0" || trimmed == "0.0" || trimmed == `"0"`
}

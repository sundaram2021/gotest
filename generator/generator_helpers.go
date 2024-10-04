package generator

import (
	"fmt"
	"go/ast"
	"math"
	"strconv"
	"strings"
)

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

package semantic

import "teslang-compiler/internal/ast"

func validType(t ast.Type) bool {
	switch t {
	case ast.TypeInt, ast.TypeVector, ast.TypeStr, ast.TypeMStr, ast.TypeBool, ast.TypeNull:
		return true
	default:
		return false
	}
}

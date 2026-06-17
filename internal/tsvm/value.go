package tsvm

import (
	"fmt"
	"strconv"
	"strings"
)

type ValueKind int

const (
	NullValue ValueKind = iota
	IntValue
	StringValue
	VectorValue
)

type Value struct {
	Kind   ValueKind
	Int    int
	String string
	Vector []Value
}

func Null() Value            { return Value{Kind: NullValue} }
func Int(n int) Value        { return Value{Kind: IntValue, Int: n} }
func String(s string) Value  { return Value{Kind: StringValue, String: s} }
func Vector(v []Value) Value { return Value{Kind: VectorValue, Vector: v} }
func Bool(ok bool) Value {
	if ok {
		return Int(1)
	}
	return Int(0)
}
func (v Value) Truthy() bool { return v.asInt() != 0 }
func (v Value) Stringify() string {
	switch v.Kind {
	case IntValue:
		return strconv.Itoa(v.Int)
	case StringValue:
		return v.String
	case VectorValue:
		parts := make([]string, len(v.Vector))
		for i, e := range v.Vector {
			parts[i] = e.Stringify()
		}
		return "[" + strings.Join(parts, ", ") + "]"
	default:
		return "null"
	}
}

func (v Value) asInt() int {
	if v.Kind == IntValue {
		return v.Int
	}
	return 0
}

func equal(a, b Value) bool {
	if a.Kind != b.Kind {
		return false
	}
	switch a.Kind {
	case IntValue:
		return a.Int == b.Int
	case StringValue:
		return a.String == b.String
	case NullValue:
		return true
	case VectorValue:
		if len(a.Vector) != len(b.Vector) {
			return false
		}
		for i := range a.Vector {
			if !equal(a.Vector[i], b.Vector[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func parseLiteral(s string) (Value, bool, error) {
	if strings.HasPrefix(s, "r") {
		return Null(), false, nil
	}
	if strings.HasPrefix(s, `"`) {
		unquoted, err := strconv.Unquote(s)
		if err != nil {
			return Null(), false, err
		}
		return String(unquoted), true, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return Null(), false, fmt.Errorf("unknown operand %q", s)
	}
	return Int(n), true, nil
}

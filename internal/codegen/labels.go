package codegen

import "fmt"

type labels struct{ n int }

func (l *labels) new(prefix string) string { l.n++; return fmt.Sprintf("%s_%d", prefix, l.n) }

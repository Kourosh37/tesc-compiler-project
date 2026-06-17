package tsvm

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type VM struct {
	program *Program
	in      *bufio.Reader
	out     io.Writer
	trace   bool
	exited  bool
	code    int
}

type Option func(*VM)

func WithInput(r io.Reader) Option {
	return func(vm *VM) {
		vm.in = bufio.NewReader(r)
	}
}

func WithOutput(w io.Writer) Option {
	return func(vm *VM) {
		vm.out = w
	}
}

func WithTrace(enabled bool) Option {
	return func(vm *VM) {
		vm.trace = enabled
	}
}

func NewVM(program *Program, opts ...Option) *VM {
	vm := &VM{program: program, in: bufio.NewReader(strings.NewReader("")), out: io.Discard}
	for _, opt := range opts {
		opt(vm)
	}
	return vm
}

func (vm *VM) Run(entry string) (int, error) {
	if entry == "" {
		entry = "main"
	}
	_, err := vm.call(entry, nil)
	return vm.code, err
}

type frame struct {
	proc *Procedure
	regs map[int]Value
	ip   int
}

func (vm *VM) call(name string, args []Value) (Value, error) {
	proc, ok := vm.program.Procedures[name]
	if !ok {
		return Null(), fmt.Errorf("procedure %q not found", name)
	}
	fr := &frame{proc: proc, regs: map[int]Value{}}
	for i, arg := range args {
		fr.regs[i+1] = arg
	}
	for fr.ip < len(proc.Instructions) {
		if vm.exited {
			return fr.regs[0], nil
		}
		inst := proc.Instructions[fr.ip]
		if vm.trace {
			fmt.Fprintf(vm.out, "# %s:%d %s %s\n", proc.Name, inst.Line, inst.Op, strings.Join(inst.Operands, ", "))
		}
		jumped, ret, err := vm.exec(fr, inst)
		if err != nil {
			return Null(), fmt.Errorf("line %d in %s: %w", inst.Line, proc.Name, err)
		}
		if ret {
			return fr.regs[0], nil
		}
		if !jumped {
			fr.ip++
		}
	}
	return fr.regs[0], nil
}

func (vm *VM) exec(fr *frame, inst Instruction) (jumped bool, returned bool, err error) {
	switch inst.Op {
	case "ret":
		return false, true, nil
	case "mov":
		if err := requireOperands(inst, 2); err != nil {
			return false, false, err
		}
		fr.set(inst.Operands[0], fr.value(inst.Operands[1]))
	case "add", "sub", "mul", "div", "mod", "lt", "gt", "le", "ge", "eq", "ne", "and", "or":
		if err := requireOperands(inst, 3); err != nil {
			return false, false, err
		}
		a, b := fr.value(inst.Operands[1]), fr.value(inst.Operands[2])
		v, err := binary(inst.Op, a, b)
		if err != nil {
			return false, false, err
		}
		fr.set(inst.Operands[0], v)
	case "not":
		if err := requireOperands(inst, 2); err != nil {
			return false, false, err
		}
		fr.set(inst.Operands[0], Bool(!fr.value(inst.Operands[1]).Truthy()))
	case "jmp":
		if err := requireOperands(inst, 1); err != nil {
			return false, false, err
		}
		ip, ok := fr.proc.Labels[inst.Operands[0]]
		if !ok {
			return false, false, fmt.Errorf("unknown label %q", inst.Operands[0])
		}
		fr.ip = ip
		return true, false, nil
	case "jz", "jnz":
		if err := requireOperands(inst, 2); err != nil {
			return false, false, err
		}
		cond := fr.value(inst.Operands[0]).Truthy()
		if (inst.Op == "jz" && !cond) || (inst.Op == "jnz" && cond) {
			ip, ok := fr.proc.Labels[inst.Operands[1]]
			if !ok {
				return false, false, fmt.Errorf("unknown label %q", inst.Operands[1])
			}
			fr.ip = ip
			return true, false, nil
		}
	case "call":
		if len(inst.Operands) < 1 {
			return false, false, fmt.Errorf("call requires a function name")
		}
		if err := vm.execCall(fr, inst.Operands); err != nil {
			return false, false, err
		}
	case "loadidx":
		if err := requireOperands(inst, 3); err != nil {
			return false, false, err
		}
		vec := fr.value(inst.Operands[1])
		idx := fr.value(inst.Operands[2]).asInt()
		if vec.Kind != VectorValue {
			return false, false, fmt.Errorf("loadidx target is not a vector")
		}
		if idx < 0 || idx >= len(vec.Vector) {
			return false, false, fmt.Errorf("vector index %d out of range", idx)
		}
		fr.set(inst.Operands[0], vec.Vector[idx])
	case "storeidx":
		if err := requireOperands(inst, 3); err != nil {
			return false, false, err
		}
		vecReg, err := registerIndex(inst.Operands[0])
		if err != nil {
			return false, false, err
		}
		vec := fr.regs[vecReg]
		idx := fr.value(inst.Operands[1]).asInt()
		if vec.Kind != VectorValue {
			return false, false, fmt.Errorf("storeidx target is not a vector")
		}
		if idx < 0 || idx >= len(vec.Vector) {
			return false, false, fmt.Errorf("vector index %d out of range", idx)
		}
		vec.Vector[idx] = fr.value(inst.Operands[2])
		fr.regs[vecReg] = vec
	default:
		return false, false, fmt.Errorf("unknown instruction %q", inst.Op)
	}
	return false, false, nil
}

func (vm *VM) execCall(fr *frame, ops []string) error {
	name := ops[0]
	switch name {
	case "read":
		if len(ops) != 2 {
			return fmt.Errorf("call read expects destination register")
		}
		text, err := vm.in.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		text = strings.TrimSpace(text)
		n, err := strconv.Atoi(text)
		if err != nil {
			return fmt.Errorf("read expected integer input")
		}
		fr.set(ops[1], Int(n))
	case "log":
		if len(ops) != 2 {
			return fmt.Errorf("call log expects one value")
		}
		fmt.Fprintln(vm.out, fr.value(ops[1]).Stringify())
	case "exit":
		if len(ops) > 2 {
			return fmt.Errorf("call exit expects at most one value")
		}
		if len(ops) == 2 {
			vm.code = fr.value(ops[1]).asInt()
		}
		vm.exited = true
	case "list":
		if len(ops) != 3 {
			return fmt.Errorf("call list expects destination and size")
		}
		size := fr.value(ops[2]).asInt()
		if size < 0 {
			return fmt.Errorf("negative list size")
		}
		items := make([]Value, size)
		for i := range items {
			items[i] = Int(0)
		}
		fr.set(ops[1], Vector(items))
	case "vector":
		if len(ops) < 2 {
			return fmt.Errorf("call vector expects destination")
		}
		items := make([]Value, 0, len(ops)-2)
		for _, op := range ops[2:] {
			items = append(items, fr.value(op))
		}
		fr.set(ops[1], Vector(items))
	case "length":
		if len(ops) != 3 {
			return fmt.Errorf("call length expects destination and vector")
		}
		vec := fr.value(ops[2])
		if vec.Kind != VectorValue {
			return fmt.Errorf("length target is not a vector")
		}
		fr.set(ops[1], Int(len(vec.Vector)))
	default:
		if len(ops) < 2 {
			return fmt.Errorf("user call requires destination register")
		}
		args := make([]Value, 0, len(ops)-2)
		for _, op := range ops[2:] {
			args = append(args, fr.value(op))
		}
		result, err := vm.call(name, args)
		if err != nil {
			return err
		}
		fr.set(ops[1], result)
	}
	return nil
}

func binary(op string, a, b Value) (Value, error) {
	switch op {
	case "add":
		if a.Kind == StringValue || b.Kind == StringValue {
			return String(a.Stringify() + b.Stringify()), nil
		}
		return Int(a.asInt() + b.asInt()), nil
	case "sub":
		return Int(a.asInt() - b.asInt()), nil
	case "mul":
		return Int(a.asInt() * b.asInt()), nil
	case "div":
		if b.asInt() == 0 {
			return Null(), fmt.Errorf("division by zero")
		}
		return Int(a.asInt() / b.asInt()), nil
	case "mod":
		if b.asInt() == 0 {
			return Null(), fmt.Errorf("modulo by zero")
		}
		return Int(a.asInt() % b.asInt()), nil
	case "lt":
		return Bool(a.asInt() < b.asInt()), nil
	case "gt":
		return Bool(a.asInt() > b.asInt()), nil
	case "le":
		return Bool(a.asInt() <= b.asInt()), nil
	case "ge":
		return Bool(a.asInt() >= b.asInt()), nil
	case "eq":
		return Bool(equal(a, b)), nil
	case "ne":
		return Bool(!equal(a, b)), nil
	case "and":
		return Bool(a.Truthy() && b.Truthy()), nil
	case "or":
		return Bool(a.Truthy() || b.Truthy()), nil
	default:
		return Null(), fmt.Errorf("unsupported binary op %q", op)
	}
}

func (fr *frame) value(op string) Value {
	if idx, err := registerIndex(op); err == nil {
		return fr.regs[idx]
	}
	v, ok, err := parseLiteral(op)
	if err != nil || !ok {
		return Null()
	}
	return v
}

func (fr *frame) set(reg string, v Value) {
	idx, err := registerIndex(reg)
	if err != nil {
		return
	}
	fr.regs[idx] = v
}

func registerIndex(s string) (int, error) {
	if !strings.HasPrefix(s, "r") {
		return 0, fmt.Errorf("expected register, got %q", s)
	}
	n, err := strconv.Atoi(strings.TrimPrefix(s, "r"))
	if err != nil || n < 0 {
		return 0, fmt.Errorf("invalid register %q", s)
	}
	return n, nil
}

func requireOperands(inst Instruction, n int) error {
	if len(inst.Operands) != n {
		return fmt.Errorf("%s expects %d operands, got %d", inst.Op, n, len(inst.Operands))
	}
	return nil
}

package spruce

import (
	"encoding/base64"
	"fmt"

	"github.com/starkandwayne/goutils/ansi"

	"github.com/starkandwayne/goutils/tree"

	log "github.com/geofffranks/spruce/log"
)

// Base64Operator ...
type Base64Operator struct{}

// Setup ...
func (Base64Operator) Setup() error {
	return nil
}

// Phase ...
func (Base64Operator) Phase() OperatorPhase {
	return EvalPhase
}

// Dependencies ...
func (Base64Operator) Dependencies(_ *Evaluator, _ []*Expr, _ []*tree.Cursor, auto []*tree.Cursor) []*tree.Cursor {
	return auto
}

// Run ...
func (Base64Operator) Run(ev *Evaluator, args []*Expr) (*Response, error) {
	log.DEBUG("running (( base64 ... )) operation at $.%s", ev.Here)
	defer log.DEBUG("done with (( base64 ... )) operation at $%s\n", ev.Here)

	if len(args) != 1 {
		return nil, fmt.Errorf("base64 operator requires exactly one string or reference argument")
	}

	var contents string

	arg := args[0]
	i := 0
	v, err := arg.Resolve(ev.Tree)
	if err != nil {
		log.DEBUG("  arg[%d]: failed to resolve expression to a concrete value", i)
		log.DEBUG("     [%d]: error was: %s", i, err)
		return nil, err
	}

	switch v.Type {
	case Literal:
		log.DEBUG("  arg[%d]: using string literal '%v'", i, v.Literal)
		log.DEBUG("     [%d]: appending '%v' to resultant string", i, v.Literal)
		if fmt.Sprintf("%T", v.Literal) != "string" {
			return nil, ansi.Errorf("@R{tried to base64 encode} @c{%v}@R{, which is not a string scalar}", v.Literal)
		}
		contents = fmt.Sprintf("%v", v.Literal)

	case Reference:
		log.DEBUG("  arg[%d]: trying to resolve reference $.%s", i, v.Reference)
		s, err := v.Reference.Resolve(ev.Tree)
		if err != nil {
			log.DEBUG("     [%d]: resolution failed\n    error: %s", i, err)
			return nil, fmt.Errorf("unable to resolve `%s`: %s", v.Reference, err)
		}

		switch s.(type) {
		case string:
			log.DEBUG("     [%d]: appending '%s' to resultant string", i, s)
			contents = fmt.Sprintf("%v", s)

		default:
			log.DEBUG("  arg[%d]: %v is not a string scalar", i, s)
			return nil, ansi.Errorf("@R{tried to base64 encode} @c{%v}@R{, which is not a string scalar}", v.Reference)
		}

	default:
		log.DEBUG("  arg[%d]: I don't know what to do with '%v'", i, arg)
		return nil, fmt.Errorf("base64 operator only accepts string literals and key reference argument")
	}
	log.DEBUG("")

	encoded := base64.StdEncoding.EncodeToString([]byte(contents))
	log.DEBUG("  resolved (( base64 ... )) operation to the string:\n    \"%s\"", string(encoded))

	return &Response{
		Type:  Replace,
		Value: string(encoded),
	}, nil
}

func init() {
	RegisterOp("base64", Base64Operator{})
}

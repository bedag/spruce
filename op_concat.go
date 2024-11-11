package spruce

import (
	"fmt"
	"strings"

	"github.com/starkandwayne/goutils/ansi"

	"github.com/starkandwayne/goutils/tree"

	log "github.com/geofffranks/spruce/log"
)

// ConcatOperator ...
type ConcatOperator struct{}

// Setup ...
func (ConcatOperator) Setup() error {
	return nil
}

// Phase ...
func (ConcatOperator) Phase() OperatorPhase {
	return EvalPhase
}

// Dependencies ...
func (ConcatOperator) Dependencies(_ *Evaluator, _ []*Expr, _ []*tree.Cursor, auto []*tree.Cursor) []*tree.Cursor {
	return auto
}

// Run ...
func (ConcatOperator) Run(ev *Evaluator, args []*Expr) (*Response, error) {
	log.DEBUG("running (( concat ... )) operation at $.%s", ev.Here)
	defer log.DEBUG("done with (( concat ... )) operation at $%s\n", ev.Here)

	var l []string

	if len(args) < 2 {
		return nil, fmt.Errorf("concat operator requires at least two arguments")
	}

	for i, arg := range args {
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
			l = append(l, fmt.Sprintf("%v", v.Literal))

		case Reference:
			log.DEBUG("  arg[%d]: trying to resolve reference $.%s", i, v.Reference)
			s, err := v.Reference.Resolve(ev.Tree)
			if err != nil {
				log.DEBUG("     [%d]: resolution failed\n    error: %s", i, err)
				return nil, fmt.Errorf("unable to resolve `%s`: %s", v.Reference, err)
			}

			switch s.(type) {
			case map[interface{}]interface{}:
				log.DEBUG("  arg[%d]: %v is not a string scalar", i, s)
				return nil, ansi.Errorf("@R{tried to concat} @c{%s}@R{, which is not a string scalar}", v.Reference)

			case []interface{}:
				log.DEBUG("  arg[%d]: %v is not a string scalar", i, s)
				return nil, ansi.Errorf("@R{tried to concat} @c{%s}@R{, which is not a string scalar}", v.Reference)

			default:
				log.DEBUG("     [%d]: appending '%s' to resultant string", i, s)
				l = append(l, fmt.Sprintf("%v", s))
			}

		default:
			log.DEBUG("  arg[%d]: I don't know what to do with '%v'", i, arg)
			return nil, fmt.Errorf("concat operator only accepts string literals and key reference arguments")
		}
		log.DEBUG("")
	}

	final := strings.Join(l, "")
	log.DEBUG("  resolved (( concat ... )) operation to the string:\n    \"%s\"", final)

	return &Response{
		Type:  Replace,
		Value: final,
	}, nil
}

func init() {
	RegisterOp("concat", ConcatOperator{})
}

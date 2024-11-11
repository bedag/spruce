package spruce

import (
	"fmt"
	"sort"

	"github.com/starkandwayne/goutils/ansi"

	log "github.com/geofffranks/spruce/log"
	"github.com/starkandwayne/goutils/tree"
)

// KeysOperator ...
type KeysOperator struct{}

// Setup ...
func (KeysOperator) Setup() error {
	return nil
}

// Phase ...
func (KeysOperator) Phase() OperatorPhase {
	return EvalPhase
}

// Dependencies ...
func (KeysOperator) Dependencies(_ *Evaluator, _ []*Expr, _ []*tree.Cursor, auto []*tree.Cursor) []*tree.Cursor {
	return auto
}

// Run ...
func (KeysOperator) Run(ev *Evaluator, args []*Expr) (*Response, error) {
	log.DEBUG("running (( keys ... )) operation at $.%s", ev.Here)
	defer log.DEBUG("done with (( keys ... )) operation at $%s\n", ev.Here)

	var vals []string

	for i, arg := range args {
		v, err := arg.Resolve(ev.Tree)
		if err != nil {
			log.DEBUG("     [%d]: resolution failed\n    error: %s", i, err)
			return nil, err
		}

		switch v.Type {
		case Literal:
			log.DEBUG("  arg[%d]: found string literal '%s'", i, v.Literal)
			log.DEBUG("           (keys operator only handles references to other parts of the YAML tree)")
			return nil, fmt.Errorf("keys operator only accepts key reference arguments")

		case Reference:
			log.DEBUG("  arg[%d]: trying to resolve reference $.%s", i, v.Reference)
			s, err := v.Reference.Resolve(ev.Tree)
			if err != nil {
				log.DEBUG("     [%d]: resolution failed\n    error: %s", i, err)
				return nil, fmt.Errorf("unable to resolve `%s`: %s", v.Reference, err)
			}

			m, ok := s.(map[interface{}]interface{})
			if !ok {
				log.DEBUG("     [%d]: resolved to something that is not a map.  that is unacceptable.", i)
				return nil, ansi.Errorf("@c{%s} @R{is not a map}", v.Reference)
			}
			log.DEBUG("     [%d]: resolved to a map; extracting keys", i)
			for k := range m {
				vals = append(vals, k.(string))
			}

		default:
			log.DEBUG("  arg[%d]: I don't know what to do with '%v'", i, arg)
			return nil, fmt.Errorf("keys operator only accepts key reference arguments")
		}
		log.DEBUG("")
	}

	switch len(args) {
	case 0:
		log.DEBUG("  no arguments supplied to (( keys ... )) operation.  oops.")
		return nil, ansi.Errorf("no arguments specified to @c{(( keys ... ))}")

	default:
		sort.Strings(vals)
		return &Response{
			Type:  Replace,
			Value: vals,
		}, nil
	}
}

func init() {
	RegisterOp("keys", KeysOperator{})
}

package spruce

import (
	"fmt"

	"github.com/starkandwayne/goutils/ansi"
	"github.com/starkandwayne/goutils/tree"

	log "github.com/geofffranks/spruce/log"
)

// GrabOperator ...
type GrabOperator struct{}

// Setup ...
func (GrabOperator) Setup() error {
	return nil
}

// Phase ...
func (GrabOperator) Phase() OperatorPhase {
	return EvalPhase
}

// Dependencies ...
func (GrabOperator) Dependencies(_ *Evaluator, _ []*Expr, _ []*tree.Cursor, auto []*tree.Cursor) []*tree.Cursor {
	return auto
}

// Run ...
func (GrabOperator) Run(ev *Evaluator, args []*Expr) (*Response, error) {
	log.DEBUG("running (( grab ... )) operation at $.%s", ev.Here)
	defer log.DEBUG("done with (( grab ... )) operation at $%s\n", ev.Here)

	var vals []interface{}

	for i, arg := range args {
		v, err := arg.Resolve(ev.Tree)
		if err != nil {
			log.DEBUG("     [%d]: resolution failed\n    error: %s", i, err)
			return nil, err
		}

		switch v.Type {
		case Literal:
			log.DEBUG("  arg[%d]: found string literal '%s'", i, v.Literal)
			vals = append(vals, v.Literal)

		case Reference:
			log.DEBUG("  arg[%d]: trying to resolve reference $.%s", i, v.Reference)
			s, err := v.Reference.Resolve(ev.Tree)
			if err != nil {
				log.DEBUG("     [%d]: resolution failed\n    error: %s", i, err)
				return nil, fmt.Errorf("unable to resolve `%s`: %s", v.Reference, err)
			}
			log.DEBUG("     [%d]: resolved to a value (could be a map, a list or a scalar); appending", i)
			vals = append(vals, s)

		default:
			log.DEBUG("  arg[%d]: I don't know what to do with '%v'", i, arg)
			return nil, fmt.Errorf("grab operator only accepts key reference arguments")
		}
		log.DEBUG("")
	}

	switch len(args) {
	case 0:
		log.DEBUG("  no arguments supplied to (( grab ... )) operation.  oops.")
		return nil, ansi.Errorf("no arguments specified to @c{(( grab ... ))}")

	case 1:
		log.DEBUG("  called with only one argument; returning value as-is")
		return &Response{
			Type:  Replace,
			Value: vals[0],
		}, nil

	default:
		log.DEBUG("  called with more than one arguments; flattening top-level lists into a single list")
		flat := []interface{}{}
		for i, lst := range vals {
			switch lst := lst.(type) {
			case []interface{}:
				log.DEBUG("    [%d]: $.%s is a list; flattening it out", i, args[i].Reference)
				flat = append(flat, lst...)
			default:
				log.DEBUG("    [%d]: $.%s is not a list; appending it as-is", i, args[i].Reference)
				flat = append(flat, lst)
			}
		}
		log.DEBUG("")

		return &Response{
			Type:  Replace,
			Value: flat,
		}, nil
	}
}

func init() {
	RegisterOp("grab", GrabOperator{})
}

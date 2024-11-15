package spruce

import (
	"github.com/starkandwayne/goutils/tree"

	log "github.com/bedag/spruce/log"
)

var keysToPrune []string

func addToPruneListIfNecessary(paths ...string) {
	for _, path := range paths {
		if !isIncluded(keysToPrune, path) {
			log.DEBUG("adding '%s' to the list of paths to prune", path)
			keysToPrune = append(keysToPrune, path)
		}
	}
}

func isIncluded(list []string, name string) bool {
	for _, entry := range list {
		if entry == name {
			return true
		}
	}

	return false
}

// PruneOperator ...
type PruneOperator struct{}

// Setup ...
func (PruneOperator) Setup() error {
	return nil
}

// Phase ...
func (PruneOperator) Phase() OperatorPhase {
	return EvalPhase
}

// Dependencies ...
func (PruneOperator) Dependencies(_ *Evaluator, _ []*Expr, _ []*tree.Cursor, auto []*tree.Cursor) []*tree.Cursor {
	return auto
}

// Run ...
func (PruneOperator) Run(ev *Evaluator, args []*Expr) (*Response, error) {
	log.DEBUG("running (( prune ... )) operation at $.%s", ev.Here)
	defer log.DEBUG("done with (( prune ... )) operation at $.%s\n", ev.Here)

	addToPruneListIfNecessary(ev.Here.String())

	// simply replace it with nil (will be pruned at the end anyway)
	return &Response{
		Type:  Replace,
		Value: nil,
	}, nil
}

func init() {
	RegisterOp("prune", PruneOperator{})
}

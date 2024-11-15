package spruce

import (
	"fmt"
	"strings"

	"github.com/starkandwayne/goutils/ansi"
	"github.com/starkandwayne/goutils/tree"

	log "github.com/bedag/spruce/log"
)

// JoinOperator is invoked with (( join <separator> <lists/strings>... )) and
// joins lists and strings into one string, separated by <separator>
type JoinOperator struct{}

// Setup ...
func (JoinOperator) Setup() error {
	return nil
}

// Phase ...
func (JoinOperator) Phase() OperatorPhase {
	return EvalPhase
}

// Dependencies returns the nodes that (( join ... )) requires to be resolved
// before its evaluation. Returns no dependencies on error, because who cares
// about eval order if Run is going to bomb out anyway.
func (JoinOperator) Dependencies(ev *Evaluator, args []*Expr, _ []*tree.Cursor, auto []*tree.Cursor) []*tree.Cursor {
	log.DEBUG("Calculating dependencies for (( join ... ))")
	deps := []*tree.Cursor{}
	if len(args) < 2 {
		log.DEBUG("Not enough arguments to (( join ... ))")
		return []*tree.Cursor{}
	}

	//skip the separator arg
	for _, arg := range args[1:] {
		if arg.Type == Literal {
			continue
		}
		if arg.Type != Reference {
			log.DEBUG("(( join ... )) argument not Literal or Reference type")
			return []*tree.Cursor{}
		}
		//get the real cursor
		finalCursor, err := arg.Resolve(ev.Tree)
		if err != nil {
			log.DEBUG("Could not resolve to a canonical path '%s'", arg.String())
			return []*tree.Cursor{}
		}
		//get the list at this location
		list, err := finalCursor.Reference.Resolve(ev.Tree)
		if err != nil {
			log.DEBUG("Could not retrieve object at path '%s'", arg.String())
			return []*tree.Cursor{}
		}
		//must be a list or a string
		switch list.(type) {
		case []interface{}:
			//add .* to the end of the cursor so we can glob all the elements
			globCursor, err := tree.ParseCursor(fmt.Sprintf("%s.*", finalCursor.Reference.String()))
			if err != nil {
				log.DEBUG("Could not parse cursor with '.*' appended. This is a BUG")
				return []*tree.Cursor{}
			}
			//have the cursor library get all the subelements for us
			subElements, err := globCursor.Glob(ev.Tree)
			if err != nil {
				log.DEBUG("Could not retrieve subelements at path '%s'. This may be a BUG.", arg.String())
				return []*tree.Cursor{}
			}
			deps = append(deps, subElements...)
		case string:
			deps = append(deps, finalCursor.Reference)
		default:
			log.DEBUG("Unsupported type at object location")
			return []*tree.Cursor{}
		}
	}

	//Append on the auto-generated deps (the operator path args)
	deps = append(deps, auto...)

	log.DEBUG("Dependencies for (( join ... )):")
	for i, dep := range deps {
		log.DEBUG("\t#%d %s", i, dep.String())
	}
	return deps
}

// Run ...
func (JoinOperator) Run(ev *Evaluator, args []*Expr) (*Response, error) {
	log.DEBUG("running (( join ... )) operation at $.%s", ev.Here)
	defer log.DEBUG("done with (( join ... )) operation at $%s\n", ev.Here)

	if len(args) == 0 {
		log.DEBUG("  no arguments supplied to (( join ... )) operation.")
		return nil, ansi.Errorf("no arguments specified to @c{(( join ... ))}")
	}

	if len(args) == 1 {
		log.DEBUG("  too few arguments supplied to (( join ... )) operation.")
		return nil, ansi.Errorf("too few arguments supplied to @c{(( join ... ))}")
	}

	var separator string
	var list []string

	for i, arg := range args {
		if i == 0 { // argument #0: separator
			sep, err := arg.Resolve(ev.Tree)
			if err != nil {
				log.DEBUG("     [%d]: resolution failed\n    error: %s", i, err)
				return nil, err
			}

			if sep.Type != Literal {
				log.DEBUG("     [%d]: unsupported type for join operator separator argument: '%v'", i, sep)
				return nil, fmt.Errorf("join operator only accepts literal argument for the separator")
			}

			log.DEBUG("     [%d]: list separator will be: %s", i, sep)
			separator = sep.Literal.(string)

		} else { // argument #1..n: list, or literal
			ref, err := arg.Resolve(ev.Tree)
			if err != nil {
				log.DEBUG("     [%d]: resolution failed\n    error: %s", i, err)
				return nil, err
			}

			switch ref.Type {
			case Literal:
				log.DEBUG("     [%d]: adding literal %s to the list", i, ref)
				list = append(list, fmt.Sprintf("%v", ref.Literal))

			case Reference:
				log.DEBUG("     [%d]: trying to resolve reference $.%s", i, ref.Reference)
				s, err := ref.Reference.Resolve(ev.Tree)
				if err != nil {
					log.DEBUG("     [%d]: resolution failed with error: %s", i, err)
					return nil, fmt.Errorf("unable to resolve `%s`: %s", ref.Reference, err)
				}

				switch s.(type) {
				case []interface{}:
					log.DEBUG("     [%d]: $.%s is a list", i, ref.Reference)
					for idx, entry := range s.([]interface{}) {
						switch entry.(type) {
						case []interface{}:
							log.DEBUG("     [%d]: entry #%d in list is a list (not a literal)", i, idx)
							return nil, ansi.Errorf("entry #%d in list is not compatible for @c{(( join ... ))}", idx)

						case map[interface{}]interface{}:
							log.DEBUG("     [%d]: entry #%d in list is a map (not a literal)", i, idx)
							return nil, ansi.Errorf("entry #%d in list is not compatible for @c{(( join ... ))}", idx)

						default:
							list = append(list, fmt.Sprintf("%v", entry))
						}
					}

				case map[interface{}]interface{}:
					log.DEBUG("     [%d]: $.%s is a map (not a list or a literal)", i, ref.Reference)
					return nil, ansi.Errorf("referenced entry is not a list or string for @c{(( join ... ))}")

				default:
					log.DEBUG("     [%d]: $.%s is a literal", i, ref.Reference)
					list = append(list, fmt.Sprintf("%v", s))
				}

			default:
				log.DEBUG("     [%d]: unsupported type for join operator: '%v'", i, ref)
				return nil, fmt.Errorf("join operator only lists with string entries, and literals as data arguments")
			}
		}
	}

	// finally, join and return
	log.DEBUG("  joined list: %s", strings.Join(list, separator))
	return &Response{
		Type:  Replace,
		Value: strings.Join(list, separator),
	}, nil
}

func init() {
	RegisterOp("join", JoinOperator{})
}

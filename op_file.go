package spruce

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/starkandwayne/goutils/ansi"

	"github.com/starkandwayne/goutils/tree"

	log "github.com/geofffranks/spruce/log"
)

// FileOperator ...
type FileOperator struct{}

// Setup ...
func (FileOperator) Setup() error {
	return nil
}

// Phase ...
func (FileOperator) Phase() OperatorPhase {
	return EvalPhase
}

// Dependencies ...
func (FileOperator) Dependencies(_ *Evaluator, _ []*Expr, _ []*tree.Cursor, auto []*tree.Cursor) []*tree.Cursor {
	return auto
}

// Run ...
func (FileOperator) Run(ev *Evaluator, args []*Expr) (*Response, error) {
	log.DEBUG("running (( file ... )) operation at $.%s", ev.Here)
	defer log.DEBUG("done with (( file ... )) operation at $%s\n", ev.Here)

	if len(args) != 1 {
		return nil, fmt.Errorf("file operator requires exactly one string or reference argument")
	}

	var fname string
	fbasepath := os.Getenv("SPRUCE_FILE_BASE_PATH")

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
		fname = fmt.Sprintf("%v", v.Literal)

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
			return nil, ansi.Errorf("@R{tried to read file} @c{%s}@R{, which is not a string scalar}", v.Reference)

		case []interface{}:
			log.DEBUG("  arg[%d]: %v is not a string scalar", i, s)
			return nil, ansi.Errorf("@R{tried to read file} @c{%s}@R{, which is not a string scalar}", v.Reference)

		default:
			log.DEBUG("     [%d]: appending '%s' to resultant string", i, s)
			fname = fmt.Sprintf("%v", s)
		}

	default:
		log.DEBUG("  arg[%d]: I don't know what to do with '%v'", i, arg)
		return nil, fmt.Errorf("file operator only accepts string literals and key reference argument")
	}
	log.DEBUG("")

	if !filepath.IsAbs(fname) {
		fname = filepath.Join(fbasepath, fname)
	}

	contents, err := os.ReadFile(fname)
	if err != nil {
		log.DEBUG("  File %s cannot be read: %s", fname, err)
		return nil, ansi.Errorf("@R{tried to read file} @c{%s}@R{: could not be read - %s}", fname, err)
	}

	log.DEBUG("  resolved (( file ... )) operation to the string:\n    \"%s\"", string(contents))

	return &Response{
		Type:  Replace,
		Value: string(contents),
	}, nil
}

func init() {
	RegisterOp("file", FileOperator{})
}

package phase2

import (
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/abtransitionit/gocore/logx"
)

// Description: returns a view of the phases in the workflow
//
// Notes:
//
// - Allow to view the workflow phases

func (wf *Workflow) GetPhaseView() (string, error) {
	var b strings.Builder

	// Header
	b.WriteString("Phase\tDescription\tTarget\tFn\tParam\n")

	// Topologically sort phases
	sorted, err := wf.topoSortByPhase()
	if err != nil {
		fmt.Println("Error sorting workflow:", err)
		return "", err
	}

	// Iterate over sorted phases
	for _, p := range sorted {
		// node
		node := p.Node
		if node == "" {
			node = "none"
		}

		// fn
		fn := p.FnAlias
		if fn == "" {
			fn = "none"
		}

		// param
		params := "none"
		if len(p.Param) > 0 {
			params = strings.Join(p.Param, ", ")
		}

		// description
		desc := p.Description
		if desc == "" {
			desc = "none"
		}

		fmt.Fprintf(&b, "%s\t%s\t%s\t%s\t%s\n", p.Name, p.Description, node, fn, params)
	}

	return b.String(), nil
}

func (wf *Workflow) GetTierView(tierList [][]Phase, logger logx.Logger) (string, error) {

	var b strings.Builder
	// sep := "-\t-\t-\t-\t-\t-\n" // ✅ separator row

	// Table header (no Params column anymore)
	// b.WriteString("Tier\tIdP\tPhase\tExe Node\tDescription\tDependencies\n")
	b.WriteString("Tier\tIdP\tPhase\tTarget\tParam\n")

	// Iterate through tiers
	for tierIndex, tierList := range tierList {
		tierID := tierIndex + 1
		for phaseIndex, p := range tierList {
			idp := phaseIndex + 1

			// deps := "none"
			// if len(p.Dependency) > 0 {
			// 	deps = strings.Join(p.Dependency, ", ")
			// }

			node := p.Node
			if node == "" {
				node = "none"
			}

			param := "none"
			if len(p.Param) > 0 { // assuming Param is a slice of strings
				param = strings.Join(p.Param, ", ")
			}

			b.WriteString(fmt.Sprintf("%d\t%d\t%s\t%s\t%s\n",
				tierID, idp, p.Name, node, param))
		}

		// b.WriteString(sep)
	}

	return b.String(), nil
}

func (wf *Workflow) GetFunctionView(cmdPathName string, registry *FnRegistry) (string, error) {

	cmdBase := filepath.Base(cmdPathName)
	keys := registry.List(cmdBase)
	if len(keys) == 0 {
		return "", nil
	}

	var b strings.Builder
	b.WriteString("Key\tModule\tFunction\n")

	for _, key := range keys {
		fn, ok := registry.Get(cmdBase, key)
		if !ok {
			continue
		}

		fnVal := reflect.ValueOf(fn)
		ptr := fnVal.Pointer()
		rf := runtime.FuncForPC(ptr)

		module, fnName := "<??>", "<??>"
		if rf != nil {
			fullName := strings.TrimPrefix(rf.Name(), "github.com/abtransitionit/")
			module = path.Dir(fullName)
			fnName = path.Base(fullName)
		}

		fmt.Fprintf(&b, "%s\t%s\t%s\n", key, module, fnName)
	}

	return b.String(), nil
}

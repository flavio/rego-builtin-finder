package rego

import (
	"fmt"
	"os"

	"github.com/deckarep/golang-set"
	"github.com/open-policy-agent/opa/ast"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/flavio/rego-builtin-finder/wasm"
)

// BuiltinsInspector inspect a Rego policy and returns the list of
// builtins used by the policy that must be provided by the runtime (SDK)
type BuiltinsInspector struct {
	// All the builtins supported by the Rego language
	all mapset.Set

	// All the builtins that have to be provided by the SDK
	sdkDependent mapset.Set
}

// NewInspector creates a new instance of BuiltinsInspector
func NewInspector() BuiltinsInspector {
	// Get all the builtins supported by Rego
	capabilities := ast.CapabilitiesForThisVersion()
	all := mapset.NewSet()
	for _, b := range capabilities.Builtins {
		all.Add(b.Name)
	}
	log.Debug().
		Interface("builtins", all).
		Msg("all rego builtins")

	// Get the builtins that are provided by Rego Wasm modules
	wasmBuiltins := wasm.Builtins()

	// Manual fixes
	// These builtins are provied by Wasm modules, despite what
	// the file we copied from OPA says
	wasmBuiltins.Add("assign")
	wasmBuiltins.Add("eq")

	// builtins that must be provided by the SDK
	sdkDependent := all.Difference(wasmBuiltins)

	return BuiltinsInspector{
		all:          all,
		sdkDependent: sdkDependent,
	}
}

// InspectPolicy loads the Rego policy from disk and then parses the AST
// tree looking for the builtins used by the policy.
//
// It returns a Set with the names of the builtins used by the policy that must
// be provided by the SDK
func (i *BuiltinsInspector) InspectPolicy(filename string) (mapset.Set, error) {
	log.Debug().Str("file", filename).Msg("inspecting")
	input, err := os.ReadFile(filename)
	if err != nil {
		return mapset.NewSet(), errors.Wrapf(err, "Error reading rego file: %s", filename)
	}

	rule, err := ast.ParseModule(filename, string(input))
	if err != nil {
		return mapset.NewSet(), errors.Wrapf(err, "Error while parsing rego file: %s", filename)
	}

	usedBuiltins := mapset.NewSet()
	ast.WalkRefs(rule, func(r ast.Ref) bool {
		log.Debug().
			Msg(fmt.Sprintf("ref: |%+v|", r))
		if i.all.Contains(r.String()) {
			usedBuiltins.Add(r.String())
		}
		return false
	})

	required := usedBuiltins.Intersect(i.sdkDependent)
	log.Debug().
		Interface("all-builtins", usedBuiltins).
		Interface("sdk-dependent", required).
		Msg("builtins")

	return required, nil
}

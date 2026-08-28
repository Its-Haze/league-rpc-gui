package types

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"testing"
)

func TestValidGameMode(t *testing.T) {
	for _, gm := range GameModes() {
		if !ValidGameMode(gm) {
			t.Errorf("ValidGameMode(%q) = false, want true", gm)
		}
	}
	if ValidGameMode("NOT_A_MODE") {
		t.Error("ValidGameMode(\"NOT_A_MODE\") = true, want false")
	}
	if ValidGameMode("") {
		t.Error("ValidGameMode(\"\") = true, want false")
	}
}

// TestGameModes_MatchesConstBlock fails if a GameMode constant is added or
func TestGameModes_MatchesConstBlock(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "queue.go", nil, 0)
	if err != nil {
		t.Fatalf("parse queue.go: %v", err)
	}

	declared := map[string]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		vs, ok := n.(*ast.ValueSpec)
		if !ok {
			return true
		}
		if id, ok := vs.Type.(*ast.Ident); !ok || id.Name != "GameMode" {
			return true
		}
		for _, v := range vs.Values {
			lit, ok := v.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				continue
			}
			s, err := strconv.Unquote(lit.Value)
			if err != nil {
				t.Fatalf("unquote %s: %v", lit.Value, err)
			}
			declared[s] = true
		}
		return true
	})

	listed := map[string]bool{}
	for _, gm := range GameModes() {
		listed[string(gm)] = true
	}

	for name := range declared {
		if !listed[name] {
			t.Errorf("GameMode %q is a declared constant but missing from GameModes()", name)
		}
	}
	for name := range listed {
		if !declared[name] {
			t.Errorf("GameModes() lists %q, which is not a declared GameMode constant", name)
		}
	}
}

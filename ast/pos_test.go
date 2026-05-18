package ast_test

import (
	"testing"

	"github.com/awalterschulze/gographviz"
	"github.com/awalterschulze/gographviz/ast"
)

// TestPos_ThreadedThroughParse verifies that lexer-emitted positions
// propagate through ast.NewID into the AST nodes for nodes, edges, and
// attributes. The DOT source below is laid out so each construct sits on
// a known line; the test asserts every Pos.Line matches.
func TestPos_ThreadedThroughParse(t *testing.T) {
	src := `digraph demo {
    alpha [color="red"]
    beta [shape="box"]
    alpha -> beta
}
`
	g, err := gographviz.ParseString(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	want := map[string]int{
		"alpha": 2,
		"beta":  3,
	}
	for _, stmt := range g.StmtList {
		switch s := stmt.(type) {
		case *ast.NodeStmt:
			id := s.NodeID.ID.Value
			if got := s.Pos.Line; got != want[id] {
				t.Errorf("NodeStmt %q: Pos.Line=%d, want %d", id, got, want[id])
			}
			if got := s.NodeID.Pos.Line; got != want[id] {
				t.Errorf("NodeID  %q: Pos.Line=%d, want %d", id, got, want[id])
			}
			if got := s.NodeID.ID.Pos.Line; got != want[id] {
				t.Errorf("ID      %q: Pos.Line=%d, want %d", id, got, want[id])
			}
			// Attrs inside the bracket should also have positions on the
			// same line as the enclosing node.
			for _, alist := range s.Attrs {
				for _, attr := range alist {
					if got := attr.Pos.Line; got != want[id] {
						t.Errorf("Attr (%s=%s): Pos.Line=%d, want %d",
							attr.Field.Value, attr.Value.Value, got, want[id])
					}
				}
			}
		case *ast.EdgeStmt:
			// The edge statement is on line 4.
			if got := s.Pos.Line; got != 4 {
				t.Errorf("EdgeStmt: Pos.Line=%d, want 4", got)
			}
			// Edge source is `alpha`, declared on line 2 but referenced
			// here from line 4 — the parser produces a NodeID for the
			// reference at line 4.
			if nid, ok := s.Source.(*ast.NodeID); ok {
				if got := nid.Pos.Line; got != 4 {
					t.Errorf("EdgeStmt.Source NodeID: Pos.Line=%d, want 4", got)
				}
			}
			for _, rh := range s.EdgeRHS {
				if got := rh.Pos.Line; got != 4 {
					t.Errorf("EdgeRH: Pos.Line=%d, want 4", got)
				}
			}
		}
	}

	// Graph-level Pos: graph ID `demo` sits on line 1.
	if g.Pos.Line != 1 {
		t.Errorf("Graph.Pos.Line=%d, want 1", g.Pos.Line)
	}
	if g.ID.Value != "demo" {
		t.Errorf("Graph.ID.Value=%q, want demo", g.ID.Value)
	}
}

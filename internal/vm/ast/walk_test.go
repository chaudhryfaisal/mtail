// Copyright 2018 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package ast_test

import (
	"testing"

	"github.com/google/mtail/internal/testutil"
	"github.com/google/mtail/internal/vm/ast"
	"github.com/google/mtail/internal/vm/parser"
	"github.com/google/mtail/internal/vm/position"
	"github.com/google/mtail/internal/vm/types"
)

type testNode struct {
}

func (t testNode) Pos() *position.Position {
	return &position.Position{}
}

func (t testNode) Type() types.Type {
	return types.None
}

type testVisitor struct {
}

func (v testVisitor) VisitBefore(n ast.Node) (ast.Visitor, ast.Node) {
	return v, n
}

func (v testVisitor) VisitAfter(n ast.Node) ast.Node {
	return n
}

func TestWalkPanicsOnUnknown(t *testing.T) {
	defer func() {
		s := recover()
		if s == nil {
			t.Errorf("No panic received")
		}
	}()
	ast.Walk(testVisitor{}, testNode{})
}

type testWalker struct {
}

func (t *testWalker) VisitBefore(n ast.Node) (ast.Visitor, ast.Node) {
	switch v := n.(type) {
	case *ast.BinaryExpr:
		if v.Op == parser.DIV {
			n = &ast.IntLit{I: 4}
		}
	}
	return t, n
}

func (t *testWalker) VisitAfter(n ast.Node) ast.Node {
	switch v := n.(type) {
	case *ast.BinaryExpr:
		if v.Op == parser.MINUS {
			n = &ast.IntLit{I: 5}
		}
	}
	return n
}

func TestAstReplacement(t *testing.T) {
	var a ast.Node

	a = &ast.BinaryExpr{Lhs: &ast.BinaryExpr{Lhs: &ast.IntLit{I: 0}, Rhs: &ast.IntLit{I: 1}, Op: parser.DIV},
		Rhs: &ast.BinaryExpr{Lhs: &ast.IntLit{I: 2}, Rhs: &ast.IntLit{I: 3}, Op: parser.MINUS},
		Op:  parser.PLUS}
	tw := &testWalker{}
	a = ast.Walk(tw, a)
	expected := &ast.BinaryExpr{Lhs: &ast.IntLit{I: 4},
		Rhs: &ast.IntLit{I: 5},
		Op:  parser.PLUS}
	diff := testutil.Diff(expected, a, testutil.IgnoreUnexported(ast.BinaryExpr{}))
	if diff != "" {
		t.Error(diff)
		s := parser.Sexp{}
		t.Log("AST:\n" + s.Dump(a))
	}
}

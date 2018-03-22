package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

// structType contains a structType node and it's name. It's a convenient
// helper type, because *ast.StructType doesn't contain the name of the struct
type structType struct {
	name string
	node *ast.StructType
}

var (
	file  = flag.String("f", "", "filename to be parsed")
	list  = flag.Bool("l", false, "list files whose formatting differs from gsorder's")
	write = flag.Bool("w", false, "write result to (source) file instead of stdout")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: gsorder [flags] [path ...]\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *file != "" {
		node, err := parse(*file)
		if err != nil {
			log.Fatal(err)
		}
		structSelection(node)
	}

	os.Exit(0)
}

// parse
func parse(file string) (ast.Node, error) {
	fset := token.NewFileSet()
	var content interface{}
	return parser.ParseFile(fset, file, content, parser.ParseComments)
}

func structSelection(node ast.Node) (int, int, error) {
	structs := collectStructs(node)
	for _, st := range structs {
		// TODO: get underlying type of field to calcacute the size
		fmt.Println(st.name)
	}

	return 0, 0, nil
}

// collectStructs collects and maps structType nodes to their positions
func collectStructs(node ast.Node) map[token.Pos]*structType {
	structs := make(map[token.Pos]*structType, 0)
	collectStructures := func(n ast.Node) bool {
		var (
			t          ast.Expr
			structName string
		)

		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Type == nil {
				return true
			}
			structName = x.Name.Name
			t = x.Type
		case *ast.CompositeLit:
			t = x.Type
		}

		x, ok := t.(*ast.StructType)
		if !ok {
			return true
		}

		structs[x.Pos()] = &structType{
			name: structName,
			node: x,
		}
		return true
	}
	ast.Inspect(node, collectStructures)
	return structs
}

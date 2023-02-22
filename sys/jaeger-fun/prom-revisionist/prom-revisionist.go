package main

import (
	"fmt"
	"log"
	"os"

	"github.com/prometheus/prometheus/promql/parser"
)

func main() {
	expr, err := parser.ParseExpr(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("before:", expr.Pretty(0))

	err = parser.Walk(Revisionist{}, expr, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("after: ", expr.Pretty(0))
}

type Revisionist struct{}

func (r Revisionist) Visit(node parser.Node, path []parser.Node) (parser.Visitor, error) {
	if node == nil && path == nil {
		return nil, nil
	}

	switch val := node.(type) {
	case *parser.VectorSelector:
		if val.Name == "calls_total" {
			val.Name = "my_calls_total"
		}

		for _, label := range val.LabelMatchers {
			if label.Name == "__name__" {
				label.Value = "my_calls_total"
			}

			if label.Name == "service_name" {
				label.Name = "service"
			}
		}
	}
	return r, nil
}

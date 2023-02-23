package main

import (
	"testing"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/require"
)

func TestRevisionist(t *testing.T) {
	testCases := []struct {
		original string
		revised  string
	}{
		{
			`  sum(rate(calls_total{service_name =~ "example-service", status_code = "STATUS_CODE_ERROR"}[5m])) by (service_name)
		     / sum(rate(calls_total{service_name =~ "example-service"}[5m])) by (service_name)`,
			`  sum(rate(my_calls_total{service =~ "example-service", status_code = "STATUS_CODE_ERROR"}[5m])) by (service)
		     / sum(rate(my_calls_total{service =~ "example-service"}[5m])) by (service)`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.original, func(t *testing.T) {
			expr, err := parser.ParseExpr(testCase.original)
			require.NoError(t, err)

			expectedExpr, err := parser.ParseExpr(testCase.revised)
			require.NoError(t, err)

			err = parser.Walk(Revisionist{}, expr, nil)
			require.NoError(t, err)

			require.Equal(t, expectedExpr.Pretty(0), expr.Pretty(0))
		})
	}
}

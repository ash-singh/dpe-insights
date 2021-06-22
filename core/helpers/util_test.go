package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSuite struct {
	firstInput  []string
	secondInput []string
	result      []string
}

func TestDifference(t *testing.T) {
	var emptyResult []string

	testSuites := []testSuite{
		{
			[]string{"sendinblue", "php-review"},
			[]string{"sendinblue", "php-review"},
			emptyResult,
		},
		{
			[]string{"account", "dpe", "go-review", "node-js-review"},
			[]string{"sendinblue", "php-review", "go-review", "node-js-review"},
			[]string{"account", "dpe"},
		},
		{
			[]string{"account", "dpe", "crm", "sre"},
			[]string{"sendinblue", "php-review", "go-review", "node-js-review"},
			[]string{"account", "dpe", "crm", "sre"},
		},
	}

	for _, item := range testSuites {
		result := Difference(item.firstInput, item.secondInput)

		if assert.ElementsMatch(t, item.result, result) == false {
			t.Errorf("Difference() with args %v, %v : Failed expected %v but got '%v'",
				item.firstInput, item.secondInput, item.result, result)
		}
	}
}

package resources

import (
	"fmt"
	"testing"
)

func TestIsResourceError(t *testing.T) {
	tests := []struct {
		err error
		is  bool
	}{
		{
			err: NewResourceError("", "", ""),
			is:  true,
		},
		{
			err: fmt.Errorf("xxx"),
			is:  false,
		},
	}

	for caseI, test := range tests {
		actual := IsResourceError(test.err)
		expect := test.is
		if actual != expect {
			t.Errorf("case %d: whether test err is ResourceError should be %v", caseI, expect)
		}
	}
}

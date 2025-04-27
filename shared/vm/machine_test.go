package vm

import (
	"testing"
)

// TODO: write more complete memory partition test that checks each const
func Test_mamoryPartitioningValidation(t *testing.T) {
	if vmTextEnd != vmMemSizeWords-1 {
		t.Fatalf(
			"invalid memory layout... vmTextEnd(%d) != vmMemSizeWords(%d)",
			vmTextEnd,
			vmMemSizeWords-1,
		)
	}
}

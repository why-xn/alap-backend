package util_test

import (
	"github.com/why-xn/alap-backend/pkg/util"
	"strings"
	"testing"
)

func TestEpcGeneration(t *testing.T) {
	dateTimeInHex := util.GetCurrentTimeInHex()
	if len(dateTimeInHex) != 8 {
		t.Errorf("invalid datetime hex length. %s", dateTimeInHex)
	}

	numHex := util.IntToHex(65535)
	if len(numHex) != 4 {
		t.Errorf("invalid number hex length")
	}

	var builder strings.Builder
	builder.WriteString(dateTimeInHex)
	builder.WriteString(numHex)
	epc := builder.String()
	t.Logf("EPC: %s", epc)
	if len(epc) != 12 {
		t.Errorf("invalid epc length")
	}
}

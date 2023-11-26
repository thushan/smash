package report

import (
	"strings"
	"testing"
)

func TestCalcTotalTime(t *testing.T) {
	var data = []struct {
		elapsedNs int64
		expected  string
	}{
		{12660100, "13ms"},
		{22592034100, "23s"},
		{60592034100, "1m1s"},
		{360592034100, "6m1s"},
		{8960592034100, "2h29m0s"},
	}

	for _, item := range data {
		actual := calcTotalTime(item.elapsedNs)
		if !strings.EqualFold(actual, item.expected) {
			t.Errorf("expected time %s, got %s", item.expected, actual)
		}
	}
}

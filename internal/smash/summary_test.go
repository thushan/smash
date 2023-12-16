package smash

import (
	"strings"
	"testing"
)

func TestCalcTotalTime(t *testing.T) {
	var data = []struct {
		expected  string
		elapsedNs int64
	}{
		{elapsedNs: 12660100, expected: "13ms"},
		{elapsedNs: 22592034100, expected: "23s"},
		{elapsedNs: 60592034100, expected: "1m1s"},
		{elapsedNs: 360592034100, expected: "6m1s"},
		{elapsedNs: 8960592034100, expected: "2h29m0s"},
	}

	for _, item := range data {
		actual := calcTotalTime(item.elapsedNs)
		if !strings.EqualFold(actual, item.expected) {
			t.Errorf("expected time %s, got %s", item.expected, actual)
		}
	}
}

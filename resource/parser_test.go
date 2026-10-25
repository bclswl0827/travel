package resource

import (
	"strings"
	"testing"
	"time"
)

const parserTestTable = `P
Depth of source [km]
∆ 0. 10.
m s m s
0.0 0 00.00 0 00.00
1.00 0.00
1.0 0 01.25 0 00.00
1.00 0.00
`

func TestNewTravelTimeMapParsesTauPTable(t *testing.T) {
	m, err := NewTravelTimeMap(parserTestTable)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if got := m[0][0]; got.Duration != 0 || got.Slowness != 1 {
		t.Fatalf("zero-time arrival was not preserved: %+v", got)
	}
	if _, found := m[0][10]; found {
		t.Fatal("all-zero no-arrival cell was preserved")
	}
	if got := m[1][0].Duration; got != 1250*time.Millisecond {
		t.Fatalf("got duration %v, want 1.25s", got)
	}
}

func TestNewTravelTimeMapRejectsInvalidNumericData(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"extra field", strings.Replace(parserTestTable, "0.0 0 00.00 0 00.00", "0.0 0 00.00 0 00.00 9", 1)},
		{"duplicate depth", strings.Replace(parserTestTable, "∆ 0. 10.", "∆ 0. 0.", 1)},
		{"NaN depth", strings.Replace(parserTestTable, "∆ 0. 10.", "∆ 0. NaN", 1)},
		{"infinite distance", strings.Replace(parserTestTable, "1.0 0 01.25", "Inf 0 01.25", 1)},
		{"NaN duration", strings.Replace(parserTestTable, "01.25", "NaN", 1)},
		{"infinite duration", strings.Replace(parserTestTable, "0 01.25", "Inf 01.25", 1)},
		{"NaN slowness", strings.Replace(parserTestTable, "1.00 0.00", "NaN 0.00", 1)},
		{"negative duration", strings.Replace(parserTestTable, "0 01.25", "-1 01.25", 1)},
		{"duration overflow", strings.Replace(parserTestTable, "0 01.25", "1e100 01.25", 1)},
		{"conflicting duplicate row", parserTestTable + "1.0 0 02.50 0 00.00\n1.00 0.00\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewTravelTimeMap(test.data); err == nil {
				t.Fatal("invalid table was accepted")
			}
		})
	}
}

func TestNewTravelTimeMapPreservesIdenticalPageOverlap(t *testing.T) {
	data := parserTestTable + "1.0 0 01.25 0 00.00\n1.00 0.00\n"
	m, err := NewTravelTimeMap(data)
	if err != nil {
		t.Fatalf("identical page overlap was rejected: %v", err)
	}
	if got := m[1][0].Duration; got != 1250*time.Millisecond {
		t.Fatalf("overlap changed duration to %v", got)
	}
}

func TestNewTravelTimeMapZeroTimeRequiresZeroDistanceAndDepth(t *testing.T) {
	data := `P
Depth of source [km]
∆ 0. 10.
m s m s
0.0 0 00.00 0 00.00
1.00 1.00
1.0 0 00.00 0 00.00
1.00 1.00
`
	m, err := NewTravelTimeMap(data)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if entry, found := m[0][0]; !found || entry.Duration != 0 {
		t.Fatal("zero-distance, zero-depth arrival was discarded")
	}
	if _, found := m[0][10]; found {
		t.Fatal("zero-time arrival at nonzero depth was preserved")
	}
	if len(m[1]) != 0 {
		t.Fatal("zero-time arrivals at nonzero distance were preserved")
	}
}

func TestNewTravelTimeMapRejectsMalformedRows(t *testing.T) {
	badDepth := `P
Depth
∆ nope
m s
0 0 1
1
`
	if _, err := NewTravelTimeMap(badDepth); err == nil {
		t.Fatal("invalid depth was accepted")
	}

	badWidth := `P
Depth
∆ 0 10
m s
0 0 1
1 1
`
	if _, err := NewTravelTimeMap(badWidth); err == nil {
		t.Fatal("short travel-time row was accepted")
	}
}

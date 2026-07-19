package resource

import (
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

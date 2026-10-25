package resource

import (
	"errors"
	"fmt"
	"maps"
	"math"
	"strconv"
	"strings"
	"time"
)

func NewTravelTimeMap(data string) (TravelTimeMap, error) {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) < 6 || (len(lines)-4)%2 != 0 {
		return nil, errors.New("invalid travel time data")
	}

	depthFields := strings.Fields(lines[2])
	if len(depthFields) < 2 {
		return nil, errors.New("travel time data has no source depths")
	}
	depths := make([]float64, len(depthFields)-1)
	depthSet := make(map[float64]struct{}, len(depths))
	for i, field := range depthFields[1:] {
		depth, err := parseTableFloat(field, "source depth")
		if err != nil {
			return nil, err
		}
		if _, exists := depthSet[depth]; exists {
			return nil, fmt.Errorf("duplicate source depth %q", field)
		}
		depthSet[depth] = struct{}{}
		depths[i] = depth
	}

	entries := make(TravelTimeMap)
	for lineIdx := 4; lineIdx < len(lines); lineIdx += 2 {
		durationFields := strings.Fields(lines[lineIdx])
		slownessFields := strings.Fields(lines[lineIdx+1])
		if len(slownessFields) == 0 || len(slownessFields) > len(depths) {
			return nil, fmt.Errorf("invalid slowness row at line %d: got %d values, maximum is %d", lineIdx+2, len(slownessFields), len(depths))
		}
		if len(durationFields) != 1+2*len(slownessFields) {
			return nil, fmt.Errorf("invalid travel-time row at line %d: got %d values, want %d", lineIdx+1, len(durationFields), 1+2*len(slownessFields))
		}

		delta, err := parseTableFloat(durationFields[0], "distance")
		if err != nil {
			return nil, err
		}
		row := make(map[float64]TravelTimeEntry)
		durationFields = durationFields[1:]

		for durationIdx := range slownessFields {
			durationMin, err := parseTableFloat(durationFields[2*durationIdx], "minute part of duration")
			if err != nil {
				return nil, err
			}
			durationSec, err := parseTableFloat(durationFields[2*durationIdx+1], "second part of duration")
			if err != nil {
				return nil, err
			}

			slowness, err := parseTableFloat(slownessFields[durationIdx], "slowness")
			if err != nil {
				return nil, err
			}
			seconds := durationMin*60 + durationSec
			if durationMin < 0 || durationSec < 0 || seconds >= float64(math.MaxInt64)/float64(time.Second) {
				return nil, fmt.Errorf("invalid travel time at distance %g, depth %g", delta, depths[durationIdx])
			}
			// Missing arrivals can retain a nonzero ray parameter in the
			// publication. Only a zero-distance, zero-depth path can take 0 s.
			if seconds == 0 && (slowness == 0 || delta != 0 || depths[durationIdx] != 0) {
				continue
			}
			row[depths[durationIdx]] = TravelTimeEntry{
				Depth:    depths[durationIdx],
				Delta:    delta,
				Slowness: slowness,
				Duration: time.Duration(math.Round(seconds * float64(time.Second))),
			}
		}
		// Page overlaps may repeat an identical row, but must not overwrite
		// a conflicting row silently.
		if previous, exists := entries[delta]; exists && !maps.Equal(previous, row) {
			return nil, fmt.Errorf("conflicting travel-time rows at distance %g", delta)
		}
		entries[delta] = row
	}

	return entries, nil
}

func parseTableFloat(field, name string) (float64, error) {
	value, err := strconv.ParseFloat(field, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s %q: %w", name, field, err)
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("%s must be finite: %q", name, field)
	}
	return value, nil
}

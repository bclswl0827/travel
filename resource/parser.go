package resource

import (
	"errors"
	"fmt"
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
	for i, field := range depthFields[1:] {
		depth, err := strconv.ParseFloat(field, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse source depth %q: %w", field, err)
		}
		depths[i] = depth
	}

	entries := make(TravelTimeMap)
	for lineIdx := 4; lineIdx < len(lines); lineIdx += 2 {
		durationFields := strings.Fields(lines[lineIdx])
		slownessFields := strings.Fields(lines[lineIdx+1])
		if len(slownessFields) > len(depths) {
			return nil, fmt.Errorf("invalid slowness row at line %d: got %d values, maximum is %d", lineIdx+2, len(slownessFields), len(depths))
		}
		if len(durationFields) < 1+2*len(slownessFields) {
			return nil, fmt.Errorf("invalid travel-time row at line %d: got %d values, want at least %d", lineIdx+1, len(durationFields), 1+2*len(slownessFields))
		}

		delta, err := strconv.ParseFloat(durationFields[0], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse distance %q: %w", durationFields[0], err)
		}
		entries[delta] = make(map[float64]TravelTimeEntry)
		durationFields = durationFields[1:]

		for durationIdx := range slownessFields {
			durationMin, err := strconv.ParseFloat(durationFields[2*durationIdx], 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse minute part of duration: %w", err)
			}
			durationSec, err := strconv.ParseFloat(durationFields[2*durationIdx+1], 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse second part of duration: %w", err)
			}

			slowness, err := strconv.ParseFloat(slownessFields[durationIdx], 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse slowness: %w", err)
			}
			// The published AK135 tables use an all-zero cell to mark a phase
			// that does not exist at this distance and source depth.
			if durationMin == 0 && durationSec == 0 && slowness == 0 {
				continue
			}
			entries[delta][depths[durationIdx]] = TravelTimeEntry{
				Depth:    depths[durationIdx],
				Delta:    delta,
				Slowness: slowness,
				Duration: time.Duration(math.Round((durationMin*60 + durationSec) * float64(time.Second))),
			}
		}
	}

	return entries, nil
}

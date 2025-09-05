package resource

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

func NewTravelTimeMap(data string) (TravelTimeMap, error) {
	lines := strings.Split(data, "\n")
	if len(lines) < 4 {
		return nil, errors.New("invalid travel time data")
	}

	depths := lo.Map(strings.Fields(lines[2])[1:], func(line string, _ int) float64 {
		depth, _ := strconv.ParseFloat(strings.TrimSpace(line), 64)
		return depth
	})

	entries := make(TravelTimeMap)
	for lineIdx := 4; lineIdx < len(lines); lineIdx += 2 {
		if len(lines[lineIdx]) == 0 {
			continue
		}

		durationFields := strings.Fields(lines[lineIdx])
		slownessFields := strings.Fields(lines[lineIdx+1])
		delta, err := strconv.ParseFloat(durationFields[0], 64)
		if err != nil {
			return nil, err
		}

		entries[delta] = make(map[float64]TravelTimeEntry)
		durationFields = durationFields[1:]

		for durationIdx := 0; durationIdx < len(durationFields)/2; durationIdx++ {
			durationMin, err := strconv.ParseFloat(durationFields[2*durationIdx], 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse minute part of duration: %w", err)
			}
			durationSec, err := strconv.ParseFloat(durationFields[2*durationIdx+1], 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse second part of duration: %w", err)
			}

			// filter out entries with zero duration
			if durationMin == 0 && durationSec == 0 {
				continue
			}

			if durationIdx < len(depths) && durationIdx < len(slownessFields) {
				slowness, err := strconv.ParseFloat(slownessFields[durationIdx], 64)
				if err != nil {
					return nil, fmt.Errorf("failed to parse slowness: %w", err)
				}
				entries[delta][depths[durationIdx]] = TravelTimeEntry{
					Depth:    depths[durationIdx],
					Delta:    delta,
					Slowness: slowness,
					Duration: time.Duration((durationMin*60 + durationSec) * float64(time.Second)),
				}
			}
		}
	}

	return entries, nil
}

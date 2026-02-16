package model

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON decodes Tivo from JSON, reconstructing each Show as the
// correct concrete type (Episode, Movie, Series) using the "kind" field.
func (t *Tivo) UnmarshalJSON(data []byte) error {
	var aux struct {
		Name    string            `json:"name"`
		Address string            `json:"address"`
		Tsn     string            `json:"tsn"`
		Shows   []json.RawMessage `json:"shows,omitempty"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Name = aux.Name
	t.Address = aux.Address
	t.Tsn = aux.Tsn
	t.Shows = make([]Show, 0, len(aux.Shows))
	for _, raw := range aux.Shows {
		var kind struct {
			Kind ShowKind `json:"kind"`
		}
		if err := json.Unmarshal(raw, &kind); err != nil {
			return fmt.Errorf("unmarshal show kind: %w", err)
		}
		switch kind.Kind {
		case ShowKindMovie:
			var m Movie
			if err := json.Unmarshal(raw, &m); err != nil {
				return fmt.Errorf("unmarshal movie: %w", err)
			}
			t.Shows = append(t.Shows, &m)
		case ShowKindSeries:
			var s Series
			if err := json.Unmarshal(raw, &s); err != nil {
				return fmt.Errorf("unmarshal series: %w", err)
			}
			t.Shows = append(t.Shows, &s)
		case ShowKindEpisode:
			var e Episode
			if err := json.Unmarshal(raw, &e); err != nil {
				return fmt.Errorf("unmarshal episode: %w", err)
			}
			t.Shows = append(t.Shows, &e)
		default:
			return fmt.Errorf("unknown show kind: %s", kind.Kind)
		}
	}
	return nil
}

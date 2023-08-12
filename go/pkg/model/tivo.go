package model

import "github.com/tartale/go/pkg/jsontime"

func (t *Tivo) Clone() (*Tivo, error) {

	originalBytes, err := jsontime.MarshalJSON(t)
	if err != nil {
		return nil, err
	}
	var newObject Tivo
	err = jsontime.UnmarshalJSON(originalBytes, &newObject)
	if err != nil {
		return nil, err
	}

	return &newObject, nil
}

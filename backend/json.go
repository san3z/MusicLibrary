package backend

import (
	"encoding/json"
)

func marshalMusic(s *Music) ([]byte, error) {
	return json.Marshal(s)
}

func unmarshalMusic(data []byte) (*Music, error) {
	var s Music
	err := json.Unmarshal(data, &s)
	return &s, err
}

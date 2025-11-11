package acp

import "encoding/json"

func mustRawJSON[T any](v T) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

package test

import "encoding/json"

func init() {
	marshalCases = append(marshalCases,
		json.RawMessage("{}"),
	)
}

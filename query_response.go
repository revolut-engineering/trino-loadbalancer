package main

import (
	"encoding/json"
)

type queryResponse struct {
	ID      string `json:"id"`
	InfoURI string `json:"infoUri"`
	NextURI string `json:"nextUri"`
	// The rest of the fields from the response are not important for the gateway
	// check https://github.com/trinodb/trino-go-client to see the complete response
}

func parse_body(data []byte) queryResponse {
	var res queryResponse
	json.Unmarshal(data, &res)

	return res
}

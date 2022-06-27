package main

import (
	"fmt"
	"testing"
)

func TestBodyMarshaling(t *testing.T) {

	data := `{
		"id":"20220624_180159_00002_e2b5c",
		"infoUri":"http://localhost:8080/ui/query.html?20220624_180159_00002_e2b5c",
		"nextUri":"http://localhost:8080/v1/statement/queued/20220624_180159_00002_e2b5c/y26b7e4479b7e72171fb4ae3b75d03bd4f0038fdf/1",
	}`
	body := parse_body([]byte(data))

	fmt.Println(data)
	fmt.Println("Body: ", body)

	want_id := "20220624_180159_00002_e2b5c"
	if body.ID != want_id {
		t.Fatalf("body.Id = %s want %s", body.ID, want_id)
	}

	want_infoUri := "http://localhost:8080/ui/query.html?20220624_180159_00002_e2b5c"
	if body.InfoURI != want_infoUri {
		t.Fatalf("body.InfoUri = %s want %s", body.InfoURI, want_infoUri)
	}

	want_nextUri := "http://localhost:8080/v1/statement/queued/20220624_180159_00002_e2b5c/y26b7e4479b7e72171fb4ae3b75d03bd4f0038fdf/1"
	if body.NextURI != want_nextUri {
		t.Fatalf("body.NextUri = %s want %s", body.NextURI, want_nextUri)
	}
}

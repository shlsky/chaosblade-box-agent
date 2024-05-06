package handler

import (
	"fmt"
	"github.com/chaosblade-io/chaos-agent/transport"
	"testing"
)

func TestSaveProtoSet(t *testing.T) {

	handler := NewScriptHandler(nil)
	r := handler.Handle(&transport.Request{
		Params: map[string]string{
			"content":     "sleep 120",
			"installPath": "/bin/bash",
			"type":        "sh",
		},
	})
	fmt.Println(r)

}

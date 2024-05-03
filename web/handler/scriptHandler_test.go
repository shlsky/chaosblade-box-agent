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
			"content":     "print('a')",
			"installPath": "/usr/local/bin/python3",
			"type":        "python",
		},
	})
	fmt.Println(r)

}

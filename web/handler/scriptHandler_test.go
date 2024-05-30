package handler

import (
	"fmt"
	"github.com/alitto/pond"
	"github.com/chaosblade-io/chaos-agent/transport"
	"testing"
	"time"
)

func TestSaveProtoSet(t *testing.T) {
	HandlerWorkerPool = pond.New(100, 1000, pond.IdleTimeout(60*time.Second))
	defer HandlerWorkerPool.StopAndWaitFor(time.Second * 60)

	handler := NewScriptHandler(nil)
	r := handler.Handle(&transport.Request{
		Params: map[string]string{
			"content":     "ps -ef | grep \"nc -l 9999\" | grep -v grep | awk '{ print $2 }'  | xargs kill -9\n sleep 1\n cat aasssa.txt\n exit 2",
			"installPath": "/bin/bash",
			"type":        "sh",
		},
	})
	fmt.Println(r)

}

func TestSaveProtoSet1(t *testing.T) {
	HandlerWorkerPool = pond.New(100, 1000, pond.IdleTimeout(60*time.Second))
	defer HandlerWorkerPool.StopAndWaitFor(time.Second * 60)

	handler := NewScriptHandler(nil)
	r := handler.Handle(&transport.Request{
		Params: map[string]string{
			"content":     "nc -l 9999",
			"installPath": "/bin/bash",
			"type":        "sh",
		},
	})
	fmt.Println(r)

}

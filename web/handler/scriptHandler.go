/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handler

import (
	"context"
	"github.com/alitto/pond"
	"github.com/chaosblade-io/chaos-agent/transport"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

type ScriptHandler struct {
	transportClient *transport.TransportClient
}

var HandlerWorkerPool *pond.WorkerPool

func NewScriptHandler(transportClient *transport.TransportClient) *ScriptHandler {
	return &ScriptHandler{
		transportClient: transportClient,
	}
}

func (ph *ScriptHandler) Handle(request *transport.Request) *transport.Response {
	logrus.Info("Receive server run script request")
	content := request.Params["content"]
	installPath := request.Params["installPath"]
	scriptType := request.Params["type"]

	fileName := uuid.New().String()
	fileName += ".sh"
	if scriptType == "python" {
		fileName += ".py"
	}

	os.WriteFile(fileName, []byte(content), 0777)

	return ExecScript(context.Background(), installPath, fileName)
}

func ExecScript(ctx context.Context, installPath, script string) *transport.Response {

	logrus.Info("async run script")
	suc := HandlerWorkerPool.TrySubmit(func() {

		defer os.Remove(script)
		// 这里需要区分windows || linux || darwin
		var cmd *exec.Cmd = exec.Command(installPath, script)
		err := cmd.Run()
		if err != nil {
			logrus.Warningf(transport.Errors[transport.CtlExecFailed], err)
		}
	})
	// 这里需要区分windows || linux || darwin
	//var cmd *exec.Cmd = exec.Command(installPath, script)
	//err := cmd.Run()
	////output, err := cmd.CombinedOutput()
	//// 2. exec uninstall command
	if !suc {
		return transport.ReturnFail(transport.CtlExecFailed, "The worker pool refused to submit")
	}

	return transport.ReturnSuccessWithResult("success")
}

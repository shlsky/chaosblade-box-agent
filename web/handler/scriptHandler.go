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
	"github.com/google/uuid"
	"os"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/chaosblade-io/chaos-agent/transport"
)

type ScriptHandler struct {
	transportClient *transport.TransportClient
}

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
	defer func() {
		os.Remove(fileName)
	}()

	os.WriteFile(fileName, []byte(content), 0777)

	return ExecScript(context.Background(), installPath, fileName)
}

func ExecScript(ctx context.Context, installPath, script string) *transport.Response {
	newCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if ctx == context.Background() {
		ctx = newCtx
	}
	// 这里需要区分windows || linux || darwin
	var cmd *exec.Cmd = exec.CommandContext(ctx, installPath, script)
	err := cmd.Start()
	//output, err := cmd.CombinedOutput()

	// 2. exec uninstall command
	if err != nil {
		logrus.Warningf(transport.Errors[transport.CtlExecFailed], err)
		return transport.ReturnFail(transport.CtlExecFailed, err.Error())
	}

	return transport.ReturnSuccessWithResult("success")
}

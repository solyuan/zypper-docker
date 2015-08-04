// Copyright (c) 2015 SUSE LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/mssola/capture"
)

func testCommand() string {
	cmd := dockerClient.(*mockClient).lastCmd
	if len(cmd) != 3 {
		return ""
	}

	// [0]: "/bin/sh", [1]: "-c", [2]: the actual command.
	// The command is basically: "zypper ref && actual command".
	return strings.TrimSpace(strings.Split(cmd[2], "&&")[1])
}

func TestListPatchesNoImageSpecified(t *testing.T) {
	setupTestExitStatus()
	dockerClient = &mockClient{}

	buffer := bytes.NewBuffer([]byte{})
	log.SetOutput(buffer)
	capture.All(func() { listPatchesCmd(testListUpdatesContext("")) })

	if testCommand() != "" {
		t.Fatalf("The command should not have been executed")
	}
	if exitInvocations != 1 {
		t.Fatalf("Expected to have exited with error")
	}
	if !strings.Contains(buffer.String(), "Error: no image name specified") {
		t.Fatal("It should've logged something\n")
	}
}

func TestListPatchesCommandFailure(t *testing.T) {
	setupTestExitStatus()
	dockerClient = &mockClient{commandFail: true}

	buffer := bytes.NewBuffer([]byte{})
	log.SetOutput(buffer)

	capture.All(func() {
		listPatchesCmd(testListUpdatesContext("opensuse:13.2"))
	})

	if testCommand() != "zypper lp" {
		t.Fatalf("Wrong command!")
	}
	if !strings.Contains(buffer.String(), "Error: Command exited with status 1") {
		t.Fatal("It should've logged something\n")
	}
	if exitInvocations != 1 {
		t.Fatalf("Expected to have exited with error")
	}
}
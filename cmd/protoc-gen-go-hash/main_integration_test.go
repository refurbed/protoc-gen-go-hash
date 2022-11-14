// Copyright (c) 2022 Refurbed GmbH. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateFile(t *testing.T) {
	workdir := t.TempDir()

	// find all the proto files in testdata
	files, err := os.ReadDir("testdata/")
	if err != nil {
		t.Fatalf("failed to read testdata dir:\n%+v", err)
	}

	var protoFilenames []string
	var goldenFilenames []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".proto") && !file.IsDir() {
			protoFilenames = append(protoFilenames, file.Name())
		}

		if strings.HasSuffix(file.Name(), ".pb.go") && !file.IsDir() {
			goldenFilenames = append(goldenFilenames, file.Name())
		}
	}

	// Compile each file, using protoc-gen-go-hash as a plugin to protoc
	for _, source := range protoFilenames {
		args := []string{
			"--plugin=protoc-gen-go-hash=../../protoc-gen-go-hash",
			"-Itestdata", "--go-hash_out=" + workdir,
			"--go-hash_opt=paths=source_relative",
		}
		args = append(args, source)
		protoc(t, args)
	}

	// Compare each generated file to the golden version.
	for _, goldenFilename := range goldenFilenames {
		goldenFilename := goldenFilename
		t.Run(strings.TrimSuffix(goldenFilename, "_hash.pb.go"), func(t *testing.T) {
			wantPath := filepath.Join("testdata", goldenFilename)
			want, err := os.ReadFile(wantPath)
			if err != nil {
				t.Fatalf("Failed to read file:\n%+v", err)
			}

			gotPath := filepath.Join(workdir, goldenFilename)
			got, err := os.ReadFile(gotPath)
			if err != nil {
				t.Fatalf("Expected file %s not generated:\n%+v", goldenFilename, err)
			}

			if !bytes.Equal(want, got) {
				t.Errorf("got:\n%swant:\n%s", got, want)
			}
		})
	}
}

func protoc(t *testing.T, args []string) {
	t.Helper()
	cmd := exec.Command("protoc")
	cmd.Args = append(cmd.Args, args...)
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		t.Log(string(out))
	}
	if err != nil {
		t.Fatalf("protoc execution failed -> args: %s:\n%+v", strings.Join(cmd.Args, " "), err)
	}
}

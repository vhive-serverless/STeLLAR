// MIT License
//
// Copyright (c) 2020 Theodor Amariucai
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package util provides support with common functionality such as reading from files, converting units,
// running bash commands.
package util

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

//ReadFile reads a file and returns the object
func ReadFile(path string) *os.File {
	log.Debugf("Reading file from `%s`", path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Could not read file: %s", err.Error())
	}
	return file
}

//FileExists checks if a file exists on disk
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//BytesToMB transforms bytes into megabytes
func BytesToMB(sizeBytes int64) float64 {
	return float64(sizeBytes) / 1024. / 1024.
}

//MBToBytes transforms megabytes into bytes
func MBToBytes(sizeMB float64) int64 {
	return int64(sizeMB) * 1024 * 1024
}

//IntegerMin returns the minimum of two integers
func IntegerMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

//RunCommandAndLog runs a command in the terminal and logs and returns the result
func RunCommandAndLog(cmd *exec.Cmd) string {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%s: %s", fmt.Sprint(err.Error()), stderr.String())
	}
	log.Debugf("Command result: %s", out.String())
	return out.String()
}
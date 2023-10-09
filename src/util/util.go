// MIT License
//
// Copyright (c) 2020 Theodor Amariucai and EASE Lab
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
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/exec"
)

// ReadFile reads a file and returns the object
func ReadFile(path string) *os.File {
	log.Debugf("Reading file from `%s`", path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Could not read file: %s", err.Error())
	}
	return file
}

// BytesToMB transforms bytes into megabytes
func BytesToMB(sizeBytes int64) float64 {
	return float64(sizeBytes) / 1024. / 1024.
}

// MBToBytes transforms megabytes into bytes
func MBToBytes(sizeMB float64) int64 {
	return int64(sizeMB) * 1024 * 1024
}

// IntegerMin returns the minimum of two integers
func IntegerMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// RunCommandAndLog runs a command in the terminal, logs the result and returns it
func RunCommandAndLog(cmd *exec.Cmd) string {
	stdoutStderr, err := cmd.CombinedOutput()
	log.Infof("Command combined output: %s\n", stdoutStderr)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(stdoutStderr)
}

func StringContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GenerateRandLowercaseLetters(length int) string {
	const lowercaseAlphabet = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = lowercaseAlphabet[rand.Intn(len(lowercaseAlphabet))]
	}
	return string(b)
}

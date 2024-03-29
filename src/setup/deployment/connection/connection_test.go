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

package connection

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetupExternalConnection(t *testing.T) {
	Initialize("www.google.com", "", apiTemplatePathFromConnectionFolder)
	require.Nil(t, Singleton.ListAPIs(), "External connection: ListAPIs() should return nil.")
	require.Nil(t, Singleton.DeployFunction, "External connection: DeployFunction should be nil.")
	require.Nil(t, Singleton.RemoveFunction, "External connection: RemoveFunction should be nil.")
	require.Nil(t, Singleton.UpdateFunction, "External connection: UpdateFunction should be nil.")
}
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

package setup

import (
	"lambda-benchmarking/client/setup/functions/connection"
	"path"
)

func getAvailableEndpoints(endpointsDirectoryPath string, config Configuration) []connection.Endpoint {
	fileProviders := []string{"vHive"}
	cloudProviders := []string{"aws"}

	if isStringInSlice(config.Provider, fileProviders) {
		endpointsFile := readFile(path.Join(endpointsDirectoryPath, config.Provider+".json"))
		return extractProviderEndpoints(endpointsFile)
	} else if isStringInSlice(config.Provider, cloudProviders) {
		return connection.Singleton.ListAPIs()
	}

	return nil
}

func removeEndpoint(s []connection.Endpoint, i int) []connection.Endpoint {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

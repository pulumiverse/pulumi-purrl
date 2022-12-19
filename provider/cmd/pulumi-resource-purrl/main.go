// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"fmt"
	"github.com/pulumiverse/pulumi-purrl/provider/pkg/provider"
	"os"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"

	"github.com/pulumiverse/pulumi-purrl/provider/pkg/version"
)

func main() {
	version := strings.TrimPrefix(version.Version, "v")

	// This method defines the provider implemented in this repository.
	purrlProvider := provider.NewProvider()

	// This method starts serving requests using the pucurl provider.
	err := p.RunProvider("purrl", version, purrlProvider)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}

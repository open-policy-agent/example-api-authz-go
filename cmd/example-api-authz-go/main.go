// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/open-policy-agent/example-api-authz-go/internal/api"
	"github.com/open-policy-agent/example-api-authz-go/internal/version"
	"github.com/open-policy-agent/opa/logging"
	"github.com/open-policy-agent/opa/sdk"
)

var configFile = flag.String("config", "", "set the OPA config file to load")
var verbose = flag.Bool("verbose", false, "enable verbose logging")
var versionFlag = flag.Bool("version", false, "print version and exit")

func main() {

	flag.Parse()

	if *configFile == "" {
		log.Fatal("Missing required --config flag")
	}

	if *versionFlag {
		fmt.Println("Version:", version.Version)
		fmt.Println("Vcs:", version.Vcs)
		os.Exit(0)
	}

	logLevel := logging.Info
	if *verbose {
		logLevel = logging.Debug
	}
	logger := logging.New()
	logger.SetLevel(logLevel)

	ctx := context.Background()

	config, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	opa, err := sdk.New(ctx, sdk.Options{Config: bytes.NewReader(config), Logger: logger})
	if err != nil {
		log.Fatal(err)
	}
	defer opa.Stop(ctx)

	if err = api.New(opa).Run(ctx); err != nil {
		log.Fatal(err)
	}

	logger.Info("Shutting down.")
}

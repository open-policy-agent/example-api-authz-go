// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package opa

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/open-policy-agent/opa/metrics"
	"github.com/open-policy-agent/opa/plugins"
	"github.com/open-policy-agent/opa/plugins/bundle"
	"github.com/open-policy-agent/opa/plugins/logs"
	"github.com/open-policy-agent/opa/plugins/status"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/server"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
	"github.com/open-policy-agent/opa/util"
)

type config struct {
	Bundle       *json.RawMessage `json:"bundle"`
	DecisionLogs *json.RawMessage `json:"decision_logs"`
	Status       *json.RawMessage `json:"status"`
}

// OPA represents an instance of the policy engine.
type OPA struct {
	decision           string
	configBytes        []byte
	manager            *plugins.Manager
	decisionLogsPlugin *logs.Plugin
}

// Config sets the configuration file to use on the OPA instance.
func Config(fileName string) func(opa *OPA) error {
	return func(opa *OPA) error {
		bs, err := ioutil.ReadFile(fileName)
		if err != nil {
			return err
		}
		opa.configBytes = bs
		return nil
	}
}

// New returns a new OPA object.
func New(opts ...func(*OPA) error) (*OPA, error) {

	opa := &OPA{}

	for _, opt := range opts {
		if err := opt(opa); err != nil {
			return nil, err
		}
	}

	store := inmem.New()

	id, err := uuid4()
	if err != nil {
		return nil, err
	}

	opa.manager, err = plugins.New(opa.configBytes, id, store)
	if err != nil {
		return nil, err
	}

	var config config

	if err := util.Unmarshal(opa.configBytes, &config); err != nil {
		return nil, err
	}

	var bundlePlugin *bundle.Plugin

	if config.Bundle != nil {
		bundlePlugin, err = bundle.New(*config.Bundle, opa.manager)
		if err != nil {
			return nil, err
		}
		opa.manager.Register(bundlePlugin)
	}

	if config.DecisionLogs != nil {
		opa.decisionLogsPlugin, err = logs.New(*config.DecisionLogs, opa.manager)
		if err != nil {
			return nil, err
		}
		opa.manager.Register(opa.decisionLogsPlugin)
	}

	if config.Status != nil {
		statusPlugin, err := status.New(*config.Status, opa.manager)
		if err != nil {
			return nil, err
		}
		opa.manager.Register(statusPlugin)
		if bundlePlugin != nil {
			bundlePlugin.Register(string("status"), statusPlugin.Update)
		}
	}

	return opa, nil
}

// Start asynchronously starts the policy engine's plugins that download
// policies, report status, etc.
func (opa *OPA) Start(ctx context.Context) error {
	return opa.manager.Start(ctx)
}

const defaultDecision = "data.system.main"

var revisionPath = storage.MustParsePath("/system/bundle/manifest/revision")

// Bool returns a boolean policy decision.
func (opa *OPA) Bool(ctx context.Context, input interface{}, opts ...func(*rego.Rego)) (bool, error) {

	m := metrics.New()
	var decisionID string
	var revision string
	var decision bool

	err := storage.Txn(ctx, opa.manager.Store, storage.TransactionParams{}, func(txn storage.Transaction) error {

		var err error

		revision, err = getRevision(ctx, opa.manager.Store, txn)
		if err != nil {
			return err
		}

		decisionID, err = uuid4()
		if err != nil {
			return err
		}

		opts = append(opts,
			rego.Metrics(m),
			rego.Query(defaultDecision),
			rego.Input(input),
			rego.Compiler(opa.manager.GetCompiler()),
			rego.Store(opa.manager.Store),
			rego.Transaction(txn))

		rs, err := rego.New(opts...).Eval(ctx)

		if err != nil {
			return err
		} else if len(rs) == 0 {
			return fmt.Errorf("undefined decision")
		} else if b, ok := rs[0].Expressions[0].Value.(bool); !ok || len(rs) > 1 {
			return fmt.Errorf("non-boolean decision")
		} else {
			decision = b
		}

		return nil
	})

	if opa.decisionLogsPlugin != nil {
		record := &server.Info{
			Revision:   revision,
			DecisionID: decisionID,
			Timestamp:  time.Now(),
			Query:      defaultDecision,
			Input:      input,
			Error:      err,
			Metrics:    m,
		}
		if err == nil {
			var x interface{} = decision
			record.Results = &x
		}
		opa.decisionLogsPlugin.Log(ctx, record)
	}

	return decision, err
}

func uuid4() (string, error) {
	bs := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, bs)
	if n != len(bs) || err != nil {
		return "", err
	}
	bs[8] = bs[8]&^0xc0 | 0x80
	bs[6] = bs[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", bs[0:4], bs[4:6], bs[6:8], bs[8:10], bs[10:]), nil
}

func getRevision(ctx context.Context, store storage.Store, txn storage.Transaction) (string, error) {
	value, err := store.Read(ctx, txn, revisionPath)
	if err != nil {
		if storage.IsNotFound(err) {
			return "", nil
		}
		return "", err
	}
	revision, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("bad revision")
	}
	return revision, nil
}

/*
Copyright 2021 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"knative.dev/pkg/signals"

	"github.com/mattmoor/grpc-example/pkg/duplex"
	"github.com/mattmoor/grpc-example/pkg/sample"
	pb "github.com/mattmoor/grpc-example/proto"
)

type envConfig struct {
	Port int `envconfig:"PORT" default:"8080" required:"true"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}

	ctx := signals.NewContext()

	d := duplex.New(env.Port)

	pb.RegisterSampleServiceServer(d.Server, sample.NewSampleServer())
	if err := d.RegisterHandler(ctx, pb.RegisterSampleServiceHandlerFromEndpoint); err != nil {
		log.Panicf("Failed to register gateway endpoint: %v", err)
	}

	if err := d.ListenAndServe(ctx); err != nil {
		log.Panicf("ListenAndServe() = %v", err)
	}

	// This will block until a signal arrives.
	<-ctx.Done()
}

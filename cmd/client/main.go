/*
Copyright 2021 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"

	pb "github.com/mattmoor/grpc-example/proto"
)

var (
	addr = flag.String("addr", "", "The URL of the service to call.")
)

func main() {
	flag.Parse()
	if *addr == "" {
		log.Fatalf("Need the address of the sample service (-addr=...)")
	}
	ctx := signals.NewContext()

	u, err := url.Parse(*addr)
	if *addr == "" {
		log.Fatalf("Malformed address -addr=%q: %v", *addr, err)
	}

	target, opts := options(u)

	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		logging.FromContext(ctx).Fatalf("Failed to connect to the datastore: %v", err)
	}
	defer conn.Close()

	client := pb.NewSampleServiceClient(conn)

	t1, err := client.Create(ctx, &pb.Thing{
		Field1: "value1",
		Field2: "value2",
	})
	if err != nil {
		logging.FromContext(ctx).Fatalf("Failed to create thing: %v", err)
	}
	logging.FromContext(ctx).Infof("Created thing: %s", t1.Id)

	t2, err := client.Update(ctx, &pb.Thing{
		Id:     t1.Id,
		Field1: "value3",
		Field2: "value4",
	})
	if err != nil {
		logging.FromContext(ctx).Fatalf("Failed to update thing: %v", err)
	}
	logging.FromContext(ctx).Infof("Updated thing: %s", t2.Id)

	tl, err := client.List(ctx, &pb.ThingFilter{
		Id: t1.Id,
	})
	if err != nil {
		logging.FromContext(ctx).Fatalf("Failed to list thing: %v", err)
	}
	for _, t := range tl.Items {
		logging.FromContext(ctx).Infof("Listed thing: id=%s, f1=%s, f2=%s", t.Id, t.Field1, t.Field2)
	}

	if _, err := client.Delete(ctx, &pb.DeleteThingRequest{Id: t1.Id}); err != nil {
		logging.FromContext(ctx).Fatalf("Failed to delete thing: %v", err)
	}
}

func options(endpoint *url.URL) (string, []grpc.DialOption) {
	switch endpoint.Scheme {
	case "http":
		port := "80"
		// Explicit port from the user signifies we should override the scheme-based defaults.
		if endpoint.Port() != "" {
			port = endpoint.Port()
		}
		return endpoint.Hostname() + ":" + port, []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithInsecure(),
		}

	case "https":
		port := "443"
		// Explicit port from the user signifies we should override the scheme-based defaults.
		if endpoint.Port() != "" {
			port = endpoint.Port()
		}
		return endpoint.Hostname() + ":" + port, []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				MinVersion: tls.VersionTLS12,
			})),
		}

	default:
		log.Fatalf("Unsupported scheme: %q", endpoint.Scheme)
		return "unreachable", nil
	}
}

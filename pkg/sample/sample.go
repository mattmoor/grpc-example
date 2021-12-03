/*
Copyright 2021 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package sample

import (
	"context"
	"sort"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/mattmoor/grpc-example/proto"
)

func NewSampleServer() pb.SampleServiceServer {
	return &sample{
		things: make(map[string]*pb.Thing),
	}
}

type sample struct {
	pb.UnimplementedSampleServiceServer

	m      sync.RWMutex
	things map[string]*pb.Thing
}

func (p *sample) Create(_ context.Context, request *pb.Thing) (*pb.Thing, error) {
	if request.Id != "" {
		return nil, status.Errorf(codes.InvalidArgument, "creates must omit id, got: %s", request.Id)
	}
	request.Id = uuid.New().String()

	p.m.Lock()
	defer p.m.Unlock()

	p.things[request.Id] = request
	return request, nil
}

func (p *sample) List(ctx context.Context, filter *pb.ThingFilter) (*pb.ThingList, error) {
	p.m.RLock()
	defer p.m.RUnlock()

	l := &pb.ThingList{}

	if filter.Id != "" {
		l.Items = make([]*pb.Thing, 0, 1)
		if t, ok := p.things[filter.Id]; ok {
			l.Items = append(l.Items, t)
		}
	} else {
		l.Items = make([]*pb.Thing, 0, len(p.things))
		for _, t := range p.things {
			l.Items = append(l.Items, t)
		}

		// Make the result stable.
		sort.Slice(l.Items, func(i, j int) bool {
			return l.Items[i].Id < l.Items[j].Id
		})
	}

	return l, nil
}

func (p *sample) Update(_ context.Context, request *pb.Thing) (*pb.Thing, error) {
	p.m.Lock()
	defer p.m.Unlock()

	if _, ok := p.things[request.Id]; !ok {
		return nil, status.Errorf(codes.NotFound, "id %s not found", request.Id)
	}

	p.m.Lock()
	defer p.m.Unlock()

	p.things[request.Id] = request
	return request, nil
}

func (p *sample) Delete(ctx context.Context, request *pb.DeleteThingRequest) (*emptypb.Empty, error) {
	p.m.Lock()
	defer p.m.Unlock()

	if _, ok := p.things[request.Id]; !ok {
		return nil, status.Errorf(codes.NotFound, "id %s not found", request.Id)
	}

	delete(p.things, request.Id)
	return &emptypb.Empty{}, nil
}

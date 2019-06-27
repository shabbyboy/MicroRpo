package handler

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/micro/go-log"

	"MicroRpo/stream-srv/mux"
	pb "MicroRpo/stream/stream-srv/proto/stream"
	"MicroRpo/stream/stream-srv/sub"
)

// Stream is a data stream
type Stream struct {
	// Mux maps stream ids to subscribers to allow stream multiplexing
	Mux *mux.Mux
	// done notifies Stream server to stop
	done chan struct{}
}

func NewStream() (*Stream, error) {
	mux, err := mux.New()
	if err != nil {
		return nil, err
	}

	done := make(chan struct{})

	return &Stream{
		Mux:  mux,
		done: done,
	}, nil
}

// Create creates new data stream.
// It returns error if the requested stream id has already been registered.
func (s *Stream) Create(ctx context.Context, req *pb.CreateRequest, resp *pb.CreateResponse) error {
	log.Logf("Received Stream.Create request with id: %s", req.Id)

	// Add new stream to stream multiplexer
	if err := s.Mux.AddStream(req.Id, 10); err != nil {
		return fmt.Errorf("Unable to create new stream: %s", err)
	}

	return nil
}

// Publish publishes data on stream
//客户端流式rpc 客户端向服务端发送数据 单向
func (s *Stream) Publish(ctx context.Context, stream pb.Stream_PublishStream) error {
	var id string
	wg := &sync.WaitGroup{}
	errCount := 0

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Logf("Stream publisher disconnected")
			break
		}

		id = msg.Id
		if err != nil {
			log.Logf("Error publishing on stream %s: %v", id, err)
			errCount++
			continue
		}

		if errCount > 5 {
			// NOTE: this is an arbitrary selected value
			log.Logf("Error threshold reached for stream: %s", id)
			break
		}

		log.Logf("Received msg on stream: %s", id)

		wg.Add(1)
		go func(msg *pb.Message) {
			defer wg.Done()
			s.Mux.Publish(msg)
		}(msg)
	}

	// wait for all the publisher goroutine to finish
	wg.Wait()

	// remove the stream from Mux
	return s.Mux.RemoveStream(id)
}
//服务端流式rpc，服务端给客户端通过stream发送数据，单向
func (s *Stream) Subscribe(ctx context.Context, req *pb.SubscribeRequest, stream pb.Stream_SubscribeStream) error {
	log.Logf("Received Stream.Subscribe request for stream: %s", req.Id)

	id := req.Id
	errCount := 0

	sub, err := sub.NewSubscriber(stream)
	if err != nil {
		return fmt.Errorf("Failed to create new subscriber for stream %s: %s", id, err)
	}

	if err := s.Mux.AddSub(id, sub); err != nil {
		return fmt.Errorf("Failed to add %v to stream: %s", sub.ID(), id)
	}

	for {
		select {
		case <-s.done:
			log.Logf("Stopping subscriber of stream: %s", id)
			// clean up is done in Stop() function
			return nil
		case err := <-sub.ErrChan():
			if err != nil {
				log.Logf("Error receiving message on stream %s: %s", id, err)
				errCount++
			}

			// NOTE: this is an arbitrary selected value
			if errCount > 5 {
				log.Logf("Error threshold reached for subscriber %s on stream: %s", sub.ID(), id)
				return s.Mux.RemSub(id, sub)
			}
		case <-sub.Done():
			// close the stream and return
			return sub.Stream().Close()
		}
	}

	return nil
}

func (s *Stream) Stop() error {
	close(s.done)
	if err := s.Mux.Stop(); err != nil {
		return fmt.Errorf("Failed to stop stream multiplexer: %s", err)
	}

	return nil
}

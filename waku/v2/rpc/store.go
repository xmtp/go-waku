package rpc

import (
	"net/http"

	"github.com/status-im/go-waku/waku/v2/node"
	"github.com/status-im/go-waku/waku/v2/protocol/pb"
	"github.com/status-im/go-waku/waku/v2/protocol/store"
)

type StoreService struct {
	node *node.WakuNode
}

// cursor       *pb.Index
// pageSize     uint64
// asc          bool

type StorePagingOptions struct {
	PageSize uint64    `json:"pageSize,omitempty"`
	Cursor   *pb.Index `json:"cursor,omitempty"`
	Forward  bool      `json:"forward,omitempty"`
}

type StoreMessagesArgs struct {
	Topic          string             `json:"pubsubTopic,omitempty"`
	ContentFilters []string           `json:"contentFilters,omitempty"`
	StartTime      float64            `json:"startTime,omitempty"`
	EndTime        float64            `json:"endTime,omitempty"`
	PagingOptions  StorePagingOptions `json:"pagingOptions,omitempty"`
}

type StoreMessagesReply struct {
	Messages   []*pb.WakuMessage  `json:"messages,omitempty"`
	PagingInfo StorePagingOptions `json:"pagingInfo,omitempty"`
	Error      string             `json:"error,omitempty"`
}

func (s *StoreService) GetV1Messages(req *http.Request, args *StoreMessagesArgs, reply *StoreMessagesReply) error {
	options := []store.HistoryRequestOption{
		store.WithAutomaticRequestId(),
		store.WithAutomaticPeerSelection(),
		store.WithPaging(args.PagingOptions.Forward, args.PagingOptions.PageSize),
		store.WithCursor(args.PagingOptions.Cursor),
	}
	res, err := s.node.Store().Query(
		req.Context(),
		store.Query{
			Topic:         args.Topic,
			ContentTopics: args.ContentFilters,
			StartTime:     args.StartTime,
			EndTime:       args.EndTime,
		},
		options...,
	)
	if err != nil {
		log.Error("Error querying messages:", err)
		reply.Error = err.Error()
		return nil
	}
	reply.Messages = res.Messages
	reply.PagingInfo = StorePagingOptions{
		PageSize: args.PagingOptions.PageSize,
		Cursor:   res.Cursor(),
		Forward:  args.PagingOptions.Forward,
	}
	return nil
}
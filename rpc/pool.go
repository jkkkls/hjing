package rpc

import "sync"

var (
	rspHeaderPool = sync.Pool{
		New: func() interface{} {
			return &RspHeader{}
		},
	}
	reqHeaderPool = sync.Pool{
		New: func() interface{} {
			return &ReqHeader{}
		},
	}
)

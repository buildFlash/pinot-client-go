package pinot

import "context"

type clientTransport interface {
	execute(ctx context.Context, brokerAddress string, query *Request) (*BrokerResponse, error)
}

package kudago_client

/*
   protoc --go_out=plugins=grpc:. *.proto
*/

import (
	"context"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"kudago/application/microservices/api_kudago/kudago_proto"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"net/http"
)

type KudagoClient struct {
	client kudago_proto.KudagoClient
	gConn  *grpc.ClientConn
	logger logger.Logger
}

func NewKudagoClient(port string, logger logger.Logger, tracer opentracing.Tracer) (*KudagoClient, error) {
	gConn, err := grpc.Dial(
		constants.Localhost+port,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(traceutils.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		return nil, err
	}

	return &KudagoClient{client: kudago_proto.NewKudagoClient(gConn), gConn: gConn, logger: logger}, nil
}

func (k *KudagoClient) AddBasic(c context.Context, num uint64, path string) (error, int) {
	input := &kudago_proto.Input{Num: num, Path: path}

	_, err := k.client.AddBasic(c, input)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (k *KudagoClient) AddToday(c context.Context) (error, int) {
	emp := &kudago_proto.Empty{}

	_, err := k.client.AddToday(c, emp)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (k *KudagoClient) Close() {
	if err := k.gConn.Close(); err != nil {
		k.logger.Warn(err)
	}
}

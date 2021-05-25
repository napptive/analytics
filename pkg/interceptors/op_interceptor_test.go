/**
 * Copyright 2021 Napptive
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package interceptors

import (
	"context"
	"fmt"
	bqprovider "github.com/napptive/analytics/pkg/provider"
	"github.com/napptive/analytics/pkg/utils"
	"github.com/napptive/grpc-ping-go"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
)

const (
	bufSize = 1024 * 1024
)

type pingAnalytics struct{}

func (p pingAnalytics) Ping(ctx context.Context, request *grpc_ping_go.PingRequest) (*grpc_ping_go.PingResponse, error) {

	// return the response
	return &grpc_ping_go.PingResponse{
		RequestNumber: request.RequestNumber,
		Data:          fmt.Sprintf("Ping [%d] received", request.RequestNumber),
	}, nil
}

var _ = ginkgo.Context("Operation interceptor", func() {

	if !utils.RunIntegrationTests("operation") {
		log.Warn().Msg("Analytics interceptor integration tests are skipped")
		return
	}

	var s *grpc.Server
	// gRPC test listener
	var lis *bufconn.Listener

	var client grpc_ping_go.PingServiceClient
	var provider bqprovider.Provider
	var proError error

	ginkgo.BeforeEach(func() {
		lis = bufconn.Listen(bufSize)

		// Create provider
		cfg := utils.GetBigQueryConfig()

		provider, proError = bqprovider.NewBigQueryProvider(*cfg)
		gomega.Expect(proError).Should(gomega.Succeed())

		s = grpc.NewServer(WithServerOpInterceptor(provider))

		handler := pingAnalytics{}
		grpc_ping_go.RegisterPingServiceServer(s, handler)
		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatal().Errs("Server exited with error: %v", []error{err})
				return
			}
		}()

		conn, err := grpc.DialContext(context.Background(), "bufnet",
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}), grpc.WithInsecure())
		gomega.Expect(err).Should(gomega.Succeed())

		client = grpc_ping_go.NewPingServiceClient(conn)

	})
	ginkgo.AfterEach(func() {
		// Flush data
		err := provider.Flush()
		gomega.Expect(err).Should(gomega.Succeed())
		s.Stop()
		lis.Close()
	})

	ginkgo.It("should be able add ", func() {

		newCtx := utils.GenerateTestFullContext()

		// Create a context with the token
		request := grpc_ping_go.PingRequest{RequestNumber: 1}
		response, err := client.Ping(newCtx, &request)
		gomega.Expect(err).Should(gomega.Succeed())
		gomega.Expect(response).ShouldNot(gomega.BeNil())
		gomega.Expect(response.RequestNumber).Should(gomega.Equal(request.RequestNumber))
		gomega.Expect(response.Data).ShouldNot(gomega.BeEmpty())

	})

})

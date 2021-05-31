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
	"github.com/napptive/analytics/pkg/entities"
	"github.com/napptive/analytics/pkg/provider"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"time"

	"context"
)

func WithServerOpInterceptor(client provider.Provider) grpc.ServerOption {
	return grpc.UnaryInterceptor(OpInterceptor(client))

}

// OpInterceptor sends the operation info call to the analytics provider
func OpInterceptor(client provider.Provider) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		userID, agent, err := ExtractDataFromContext(ctx)
		if err != nil {
			return nil, err
		}

		// Send the data to bigQuery
		if err = client.Send(entities.Operation{
			Timestamp: time.Now(),
			UserID:    userID,
			Operation: info.FullMethod,
			Source:    agent,
		}); err != nil {
			log.Err(err).Msg("error sending analytics data")
		}

		return handler(ctx, req)
	}
}

func WithServerOpStreamInterceptor(client provider.Provider) grpc.ServerOption {
	return grpc.StreamInterceptor(OpStreamInterceptor(client))

}

// OpStreamInterceptor verifies the JWT token and adds the claim information in the context
func OpStreamInterceptor(client provider.Provider) grpc.StreamServerInterceptor {
	return func(srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {

		ctx := stream.Context()
		userID, agent, err := ExtractDataFromContext(ctx)
		if err != nil {
			return err
		}

		// Send the data to bigQuery
		if err = client.Send(entities.Operation{
			Timestamp: time.Now(),
			UserID:    userID,
			Operation: info.FullMethod,
			Source:    agent,
		}); err != nil {
			log.Err(err).Msg("error sending analytics data")
		}

		return handler(srv, stream)
	}
}
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
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"time"
)

// TimeInterceptor is an interceptor that allows you to measure the time it takes to make a call to a grpc call
// How to use:
// gRPCServer := grpc.NewServer(interceptor.WithServerTimeInterceptor())

func WithServerTimeInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(TimeInterceptor())

}

func TimeInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		start := time.Now()

		response, err := handler(ctx, req)

		log.Info().Dur("elapsed", time.Since(start)).Str("function", info.FullMethod).Msg("[TimeInterceptor] time")

		return response, err
	}
}

func WithServerTimeStreamInterceptor() grpc.ServerOption {
	return grpc.StreamInterceptor(TimeStreamInterceptor())

}

func TimeStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {

		start := time.Now()

		err := handler(srv, stream)

		log.Info().Str("elapsed", time.Since(start).String()).Str("function", info.FullMethod).Msg("[TimeInterceptor] time")

		return err
	}
}
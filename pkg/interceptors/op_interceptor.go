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
	"github.com/napptive/analytics/pkg/analytics"
	"github.com/napptive/analytics/pkg/entities"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"time"

	"context"
)

func WithServerOpInterceptor(client analytics.Provider) grpc.ServerOption {
	return grpc.UnaryInterceptor(OpInterceptor(client))

}

// OpInterceptor sends the operation info call to the analytics provider
func OpInterceptor(client analytics.Provider) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		userID, err := GetUserIDFromContext(ctx)
		if err != nil {
			return nil, err
		}

		// Send the data to bigQuery
		if err = client.SendOperationData(entities.OperationData{
			Timestamp: time.Now(),
			UserID:    userID,
			Operation: info.FullMethod,
		}); err != nil {
			log.Err(err).Msg("error sending analytics data")
		}

		return handler(ctx, req)
	}
}

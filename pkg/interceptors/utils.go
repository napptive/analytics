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
	"github.com/napptive/nerrors/pkg/nerrors"
	"google.golang.org/grpc/metadata"
)

const (
	// UserIdKey with the name of the key that will be injected in the context metadata corresponding to the user identifier.
	UserIDKey = "user_id"
	// AgentHeader with the key name for the agent payload.
 	AgentHeader = "agent"
)

// ExtractDataFromContext returns the userID and the agent
// the userID is required, but the agent is not
func ExtractDataFromContext(ctx context.Context) (string, string, error) {
	// get the userId from the context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", nerrors.NewInternalError("no metadata found").ToGRPC()
	}
	userID, exists := md[UserIDKey]
	if !exists {
		return "", "", nerrors.NewInternalError("userId not found in metadata")
	}
	agentValue := ""
	agent, exists := md[AgentHeader]
	if exists {
		agentValue = agent[0]
	}
	return userID[0], agentValue, nil
}

/**
 * Copyright 2020 Napptive
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
package entities

import "time"

// LoginData with the login info stored
type LoginData struct {
	// Timestamp with the event timestamp
	Timestamp time.Time
	// UserID with the user identifier
	UserID string
	// Method with the login source (PAT, CLI, Web)
	Method string
}

// OperationData with the operation info stored
type OperationData struct {
	// Timestamp with the event timestamp
	Timestamp time.Time
	// UserID with the user identifier
	UserID string
	// Operation with the name of the GRPC method
	Operation string
}

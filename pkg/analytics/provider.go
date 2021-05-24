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
package analytics

import (
	"github.com/napptive/analytics/pkg/entities"
)

// Provider with an interface that defines the monitoring provider methods
type Provider interface {
	// SendLoginData puts a login in the database
	SendLoginData(data entities.LoginData) error
	// SendOperationData puts an operation data in the database
	SendOperationData(data entities.OperationData) error
}


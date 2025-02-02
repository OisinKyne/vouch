// Copyright © 2021 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/rs/zerolog"
)

func init() {
	// We seed math.rand here so that we can obtain different IDs for requests.
	// This is purely used as a way to match request and response entries in logs, so there is no
	// requirement for this to cryptographically secure.
	rand.Seed(time.Now().UnixNano())
}

// LogWithID returns a new logger based on the supplied logger with an additional ID field.
func LogWithID(_ context.Context, log zerolog.Logger, tag string) zerolog.Logger {
	// #nosec G404
	return log.With().Str(tag, fmt.Sprintf("%02x", rand.Int31())).Logger()
}

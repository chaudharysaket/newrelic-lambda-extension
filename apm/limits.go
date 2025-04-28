// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apm

import "time"

const (
	// app behavior

	// collectorTimeout is the timeout used in the client for communication
	// with New Relic's servers.
	collectorTimeout = 20 * time.Second
)

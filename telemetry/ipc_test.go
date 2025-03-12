package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTelemetryChannel(t *testing.T) {
	channel,_, err := InitTelemetryChannel()

	assert.Nil(t, err)
	assert.Empty(t, channel)
}

package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RequestContext(t *testing.T) {

	requestID := "requestID"

	ctx := WithRequestCtx(context.Background(), requestID)

	actualGC, ok := RequestCtx(ctx)
	assert.True(t, ok)
	assert.Equal(t, requestID, actualGC.requestID)
}

package srv

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func assertTraceID(t *testing.T, header *http.Header, u *url.URL, expected string) {
	t.Helper()

	req := &http.Request{}
	if header != nil {
		req.Header = *header
	}

	if u != nil {
		req.URL = u
	}

	provider := newRequestTraceIDProvider(req)

	actual := provider.GetTraceID()

	assert.Equal(t, expected, actual)
}

func TestRequestTraceIDProviderFromQuery(t *testing.T) {
	u := &url.URL{RawQuery: "traceId=test-id"}

	assertTraceID(t, nil, u, "test-id")
}

func TestRequestTraceIDProviderFromHeader(t *testing.T) {
	header := http.Header{}
	header.Add("X-Trace-Id", "test-id")

	assertTraceID(t, &header, nil, "test-id")
}

func TestRequestTraceIDProviderGenerateTraceID(t *testing.T) {
	provider := newRequestTraceIDProvider(new(http.Request))

	actual := provider.GetTraceID()

	_, err := uuid.Parse(actual)

	assert.NoError(t, err)
}

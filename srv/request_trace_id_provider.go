package srv

import (
	"net/http"

	"github.com/google/uuid"
)

type requestTraceIDProvider struct {
	traceID string
	req     *http.Request
}

func (p *requestTraceIDProvider) GetTraceID() string {
	if len(p.traceID) == 0 {
		if p.req.URL != nil {
			// Try to get the trace id from the query first.
			p.traceID = p.req.URL.Query().Get(TraceIDName)
		}

		if len(p.traceID) == 0 {
			// Try to get the trace id from the headers.
			p.traceID = p.req.Header.Get(TraceIDReqHeaderName)

			// Generate trace id if there is no such in the query or in the headers.
			if len(p.traceID) == 0 {
				p.traceID = uuid.New().String()
			}
		}
	}

	return p.traceID
}

func newRequestTraceIDProvider(req *http.Request) TraceIDProvider {
	provider := new(requestTraceIDProvider)
	provider.req = req

	return provider
}

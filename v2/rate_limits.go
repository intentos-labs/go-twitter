package twitter

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

const (
	rateLimit     = "x-rate-limit-limit"
	rateRemaining = "x-rate-limit-remaining"
	rateReset     = "x-rate-limit-reset"

	rateLimit24H     = "X-User-Limit-24hour-Limit"
	rateRemaining24H = "X-User-Limit-24hour-Remaining"
	rateReset24H     = "X-User-Limit-24hour-Reset"
)

// Epoch is the UNIX seconds from 1/1/1970
type Epoch int

// Time returns the epoch time
func (e Epoch) Time() time.Time {
	return time.Unix(int64(e), 0)
}

// RateLimit are the rate limit values from the response header
type RateLimit struct {
	Limit     int
	Remaining int
	Reset     Epoch
}

func rateFromHeader24H(header http.Header) *RateLimit {
	limit, err := strconv.Atoi(header.Get(rateLimit24H))
	if err != nil {
		return nil
	}
	remaining, err := strconv.Atoi(header.Get(rateRemaining24H))
	if err != nil {
		return nil
	}
	reset, err := strconv.Atoi(header.Get(rateReset24H))
	if err != nil {
		return nil
	}
	return &RateLimit{
		Limit:     limit,
		Remaining: remaining,
		Reset:     Epoch(reset),
	}
}

func rateFromHeader(header http.Header) *RateLimit {
	limit, err := strconv.Atoi(header.Get(rateLimit))
	if err != nil {
		return nil
	}
	remaining, err := strconv.Atoi(header.Get(rateRemaining))
	if err != nil {
		return nil
	}
	reset, err := strconv.Atoi(header.Get(rateReset))
	if err != nil {
		return nil
	}
	return &RateLimit{
		Limit:     limit,
		Remaining: remaining,
		Reset:     Epoch(reset),
	}
}

// RateLimitFromError returns the rate limits from an error.  If there are not any limits, false is returned.
func RateLimitFromError(err error) (*RateLimit, bool) {
	var er *ErrorResponse
	var hr *HTTPError
	var rde *ResponseDecodeError
	switch {
	case errors.As(err, &er) && er.RateLimit != nil:
		return er.RateLimit, true
	case errors.As(err, &hr) && hr.RateLimit != nil:
		return hr.RateLimit, true
	case errors.As(err, &rde) && rde.RateLimit != nil:
		return rde.RateLimit, true
	default:
	}
	return nil, false
}

package util

import (
	"net/http"
	"strconv"
	"time"
)

func AddGetCurrentCrowdQueries(req *http.Request, since time.Time) {
	q := req.URL.Query()
	s := strconv.FormatInt(since.Unix(), 10)
	q.Add("since", s)
	req.URL.RawQuery = q.Encode()
}

func AddGetCrowdBetweenDatesQueries(req *http.Request, from, until time.Time, forEvery time.Duration) {
	q := req.URL.Query()
	fromString := strconv.FormatInt(from.Unix(), 10)
	untilString := strconv.FormatInt(until.Unix(), 10)
	forEveryString := strconv.FormatInt(int64(forEvery.Seconds()), 10)
	q.Add("from", fromString)
	q.Add("until", untilString)
	q.Add("for", forEveryString)
	req.URL.RawQuery = q.Encode()
}

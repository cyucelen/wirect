package util

import (
	"net/http"
	"strconv"
	"time"
)

func AddGetCrowdRequestHeaders(req *http.Request, since time.Time) {
	q := req.URL.Query()
	s := strconv.FormatInt(since.Unix(), 10)
	q.Add("since", s)
	req.URL.RawQuery = q.Encode()
}

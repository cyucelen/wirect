package util

import (
	"net/http"
	"strconv"
	"time"
)

func AddGetCrowdRequestHeaders(req *http.Request, since time.Time, snifferMAC string) {
	q := req.URL.Query()
	s := strconv.FormatInt(since.Unix(), 10)
	q.Add("since", s)
	q.Add("sniffer", snifferMAC)
	req.URL.RawQuery = q.Encode()
}

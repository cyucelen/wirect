package util

import (
	"net/http"
	"time"
)

func AddGetCrowdRequestHeaders(req *http.Request, since time.Time, snifferMAC string) {
	q := req.URL.Query()
	q.Add("since", since.Format(time.RFC3339))
	q.Add("sniffer", snifferMAC)
	req.URL.RawQuery = q.Encode()
}

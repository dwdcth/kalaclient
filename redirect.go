package kala

import (
	"fmt"
	"gopkg.in/resty.v1"
	"net"
	"net/http"
	"strings"
)

var (
	hdrUserAgentValue = "go-resty v%s - https://github.com/go-resty/resty"
	hdrUserAgentKey   = http.CanonicalHeaderKey("User-Agent")
)

func FlexibleRedirectPolicy(noOfRedirect int) resty.RedirectPolicy {
	return resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		if len(via) >= noOfRedirect {
			return fmt.Errorf("Stopped after %d redirects", noOfRedirect)
		}

		checkHostAndAddHeaders(req, via[0])

		return nil
	})
}

func getHostname(host string) (hostname string) {
	if strings.Index(host, ":") > 0 {
		host, _, _ = net.SplitHostPort(host)
	}
	hostname = strings.ToLower(host)
	return
}

func checkHostAndAddHeaders(cur *http.Request, pre *http.Request) {
	curHostname := getHostname(cur.URL.Host)
	preHostname := getHostname(pre.URL.Host)
	cur.Method = pre.Method
	if strings.EqualFold(curHostname, preHostname) {
		for key, val := range pre.Header {
			cur.Header[key] = val
		}
	} else { // only library User-Agent header is added
		cur.Header.Set(hdrUserAgentKey, fmt.Sprintf(hdrUserAgentValue, "1.0"))
	}
}

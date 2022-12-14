package logger

import (
	"github.com/felixge/httpsnoop"
	"github.com/kjk/siser"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func LogRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ri := &HTTPReqInfo{
			method:    r.Method,
			uri:       r.URL.String(),
			referer:   r.Header.Get("Referer"),
			userAgent: r.Header.Get("User-Agent"),
		}

		ri.ipaddr = requestGetRemoteAddress(r)

		// this runs handler h and captures information about
		// HTTP request
		m := httpsnoop.CaptureMetrics(h, w, r)

		ri.code = m.Code
		ri.size = m.Written
		ri.duration = m.Duration
		logHTTPReq(ri)
	}
	return http.HandlerFunc(fn)
}

// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]"
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// requestGetRemoteAddress returns ip address of the client making the request,
// taking into account http proxies
func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrRealIP
}

var (
	muLogHTTP sync.Mutex
)

func logHTTPReq(ri *HTTPReqInfo) {
	var rec siser.Record
	rec.Name = "httplog"
	rec.Write("method", ri.method)
	rec.Write("uri", ri.uri)
	if ri.referer != "" {
		rec.Write("referer", ri.referer)
	}
	rec.Write("ipaddr", ri.ipaddr)
	rec.Write("code", strconv.Itoa(ri.code))
	rec.Write("size", strconv.FormatInt(ri.size, 10))
	durMs := ri.duration / time.Millisecond
	rec.Write("duration", strconv.FormatInt(int64(durMs), 10))
	rec.Write("ua", ri.userAgent)

	muLogHTTP.Lock()
	defer muLogHTTP.Unlock()
	f, err := os.Create("http_access.log")
	log.Println(err)
	w := siser.NewWriter(f)
	w.WriteRecord(&rec)
	//_, _ = httpLogSiser.WriteRecord(&rec)
}

// LogReqInfo describes info about HTTP request
type HTTPReqInfo struct {
	// GET etc.
	method  string
	uri     string
	referer string
	ipaddr  string
	// response code, like 200, 404
	code int
	// number of bytes of the response sent
	size int64
	// how long did it take to
	duration  time.Duration
	userAgent string
}

package heff

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// DefaultHoneypot is an http.HandlerFunc that serves random HTML from the
// DefaultMarkovMap, 100KB at a time.
var DefaultHoneypot = NewHoneypot(DefaultMarkovMap, 100*1<<10)

// NewHoneypot creates an http.HandlerFunc from a MarkovMap
func NewHoneypot(mm MarkovMap, buffsize int) http.HandlerFunc {
	var pool sync.Pool

	getBuffer := func() []byte {
		x := pool.Get()
		if buf, ok := x.([]byte); ok {
			return buf
		}
		return make([]byte, buffsize)
	}

	putBuffer := func(buf []byte) {
		pool.Put(buf)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		s := time.Now()
		log.Printf("Start FOR: %v USER AGENT: %q PATH: %q ",
			r.Header["X-Forwarded-For"], r.UserAgent(), r.URL)
		buf := getBuffer()
		defer putBuffer(buf)
		io.WriteString(w, "<HTML>\n<BODY>\n")
		n, err := io.CopyBuffer(w, mm, buf)
		log.Printf("Finish FOR: %v USER AGENT: %q PATH: %q "+
			"BYTES: %d DURATION: %v ERROR: %v",
			r.Header["X-Forwarded-For"], r.UserAgent(), r.URL,
			n, time.Since(s), err)
	}
}

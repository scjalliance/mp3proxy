package mp3proxy

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Mp3Proxy serves an MP3 file verbatim through a proxy, but replaces the Content-Type header with `audio/mp3`
func Mp3Proxy(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	u := q.Get("mp3")
	if u == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid request.")
		return
	}
	c := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := c.Get(u)
	for hK, hV := range resp.Header {
		for _, v := range hV {
			w.Header().Add(hK, v)
		}
	}
	w.Header().Del("Content-Type")
	if err != nil {
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "audio/mp3")
	w.WriteHeader(resp.StatusCode)
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

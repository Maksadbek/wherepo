package route

import (
	"github.com/Maksadbek/wherepo/conf"
	"github.com/Maksadbek/wherepo/models"
	"github.com/Maksadbek/wherepo/logger"
	"github.com/Maksadbek/wherepo/cache"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	tollboothConf "github.com/didip/tollbooth/config"
	"github.com/garyburd/redigo/redis"
)

var config conf.App

// setup rate limiter, 1 request per 5 seconds
var Limiter *tollboothConf.Limiter

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	sniffDone bool
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.sniffDone {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", http.DetectContentType(b))
		}
		w.sniffDone = true
	}
	return w.Writer.Write(b)
}

// Wrap a http.Handler to support transparent gzip encoding.
func GzipHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Accept-Encoding")
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		h.ServeHTTP(&gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}
func Initialize(c conf.App) error {
	config = c
	err := cache.Initialize(config)
	if err != nil {
		return err
	}
	// setup rate limiter
	Limiter = tollbooth.NewLimiter(1, time.Second*5)
	Limiter.Methods = append(Limiter.Methods, "POST")
	//Limiter.IPLookups = []string{""}
	return err
}

// GetTrackers can be used to get list of trackers
// if user does not exist in cache then in caches from mysql
func GetTrackers(name string) (trackers cache.Usr, err error) {
	trackers, err = cache.UsrTrackers(name)
	logger.FuncLog("route.GetTracker", "GetTracker", nil, nil)
	if err == nil || err != redis.ErrNil {
		return
	}
	// if redis result is nil
	trackers, err = models.UsrTrackers(name)
	if err != nil {
		return
	}
	err = cache.SetUsrTrackers(trackers)
	if err != nil {
		logger.FuncLog("route.GetTrackers", "GetTrackers", nil, err)
		return
	}

	return
}

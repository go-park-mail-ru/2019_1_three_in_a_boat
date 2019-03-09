package handlers

import (
	"fmt"
	"github.com/google/logger"
	"net/http"
)

func fmtHttpLog(msg string, r *http.Request) string {
	return fmt.Sprintf("%s %s: %s %s\n",
		r.RemoteAddr, r.Method, r.RequestURI, msg)
}

func LogInfo(depth int, msg string, r *http.Request) {
	logger.InfoDepth(1+depth, fmtHttpLog(msg, r))
}

func LogWarning(depth int, msg string, r *http.Request) {
	logger.WarningDepth(1+depth, fmtHttpLog(msg, r))
}

func LogError(depth int, msg string, r *http.Request) {
	logger.ErrorDepth(1+depth, fmtHttpLog(msg, r))
}

func LogFatal(depth int, msg string, r *http.Request) {
	// it was called fatal depth hence the mediocre joke, 4/10 if I'm generous:
	// fatal depth sounds like the kind of movie you'd watch just for the sake of
	// killing 2 hours and forget about it as you're walking out of the cinema
	logger.FatalDepth(1+depth, fmt.Sprintf("%s\t%s:%s %s\n",
		r.RemoteAddr, r.Method, r.RequestURI, msg))
}

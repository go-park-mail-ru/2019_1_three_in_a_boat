package handlers

import (
	"fmt"
	"github.com/google/logger"
	"net/http"
)

func fmtHttpLog(msg string, r *http.Request) string {
	return fmt.Sprintf("%s %s:\t%s %s\n",
		r.RemoteAddr, r.Method, r.RequestURI, msg)
}

func logInfo(depth int, msg string, r *http.Request) {
	logger.InfoDepth(1+depth, fmtHttpLog(msg, r))
}

func logWarning(depth int, msg string, r *http.Request) {
	logger.WarningDepth(1+depth, fmtHttpLog(msg, r))
}

func logError(depth int, msg string, r *http.Request) {
	logger.ErrorDepth(1+depth, fmtHttpLog(msg, r))
}

func logFatal(depth int, msg string, r *http.Request) {
	// fatal depth sounds like the kind of movie you'd watch just for the sake of
	// killing 2 hours and forget about it as you're walking out of the cinema
	logger.FatalDepth(1+depth, fmt.Sprintf("%s\t%s:%s %s\n",
		r.RemoteAddr, r.Method, r.RequestURI, msg))
}

package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/logger"
)

// utility function: adds some remoteAddr, method, uri in front of the message
func fmtHttpLog(msg string, r *http.Request) string {
	return fmt.Sprintf("%s %s: %s %s\n",
		r.RemoteAddr, r.Method, r.RequestURI, msg)
}

// Logs info with some http-related metadata
func LogInfo(depth int, msg string, r *http.Request) {
	logger.InfoDepth(1+depth, fmtHttpLog(msg, r))
}

// Logs warning with some http-related metadata
func LogWarning(depth int, msg string, r *http.Request) {
	logger.WarningDepth(1+depth, fmtHttpLog(msg, r))
}

// Logs error with some http-related metadata
func LogError(depth int, msg string, r *http.Request) {
	logger.ErrorDepth(1+depth, fmtHttpLog(msg, r))
}

// Logs fatal with some http-related metadata. Terminates the program with
// code 1. Should never be used in handlers. Why is it even here?
func LogFatal(depth int, msg string, r *http.Request) {
	// it was called fatal depth hence the mediocre joke, 4/10 if I'm generous:
	// fatal depth sounds like the kind of movie you'd watch just for the sake of
	// killing 2 hours and forget about it as you're walking out of the cinema
	// I'll be honest with you this function is still here just because of the joke
	logger.FatalDepth(1+depth, fmt.Sprintf("%s\t%s:%s %s\n",
		r.RemoteAddr, r.Method, r.RequestURI, msg))
}

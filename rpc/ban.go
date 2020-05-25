package rpc

import (
	"math/rand"
	"net"
	"net/http"
)

var allHTTPStatusCodes = []int{

	http.StatusContinue,           // RFC 7231, 6.2.1
	http.StatusSwitchingProtocols, // RFC 7231, 6.2.2
	http.StatusProcessing,         // RFC 2518, 10.1
	http.StatusEarlyHints,         // RFC 8297

	http.StatusOK,                   // RFC 7231, 6.3.1
	http.StatusCreated,              // RFC 7231, 6.3.2
	http.StatusAccepted,             // RFC 7231, 6.3.3
	http.StatusNonAuthoritativeInfo, // RFC 7231, 6.3.4
	http.StatusNoContent,            // RFC 7231, 6.3.5
	http.StatusResetContent,         // RFC 7231, 6.3.6
	http.StatusPartialContent,       // RFC 7233, 4.1
	http.StatusMultiStatus,          // RFC 4918, 11.1
	http.StatusAlreadyReported,      // RFC 5842, 7.1
	http.StatusIMUsed,               // RFC 3229, 10.4.1

	http.StatusMultipleChoices,  // RFC 7231, 6.4.1
	http.StatusMovedPermanently, // RFC 7231, 6.4.2
	http.StatusFound,            // RFC 7231, 6.4.3
	http.StatusSeeOther,         // RFC 7231, 6.4.4
	http.StatusNotModified,      // RFC 7232, 4.1
	http.StatusUseProxy,         // RFC 7231, 6.4.5
	// _                      ,  // RFC 7231, 6.4.6 (Unused)
	http.StatusTemporaryRedirect, // RFC 7231, 6.4.7
	http.StatusPermanentRedirect, // RFC 7538, 3

	http.StatusBadRequest,                   // RFC 7231, 6.5.1
	http.StatusUnauthorized,                 // RFC 7235, 3.1
	http.StatusPaymentRequired,              // RFC 7231, 6.5.2
	http.StatusForbidden,                    // RFC 7231, 6.5.3
	http.StatusNotFound,                     // RFC 7231, 6.5.4
	http.StatusMethodNotAllowed,             // RFC 7231, 6.5.5
	http.StatusNotAcceptable,                // RFC 7231, 6.5.6
	http.StatusProxyAuthRequired,            // RFC 7235, 3.2
	http.StatusRequestTimeout,               // RFC 7231, 6.5.7
	http.StatusConflict,                     // RFC 7231, 6.5.8
	http.StatusGone,                         // RFC 7231, 6.5.9
	http.StatusLengthRequired,               // RFC 7231, 6.5.10
	http.StatusPreconditionFailed,           // RFC 7232, 4.2
	http.StatusRequestEntityTooLarge,        // RFC 7231, 6.5.11
	http.StatusRequestURITooLong,            // RFC 7231, 6.5.12
	http.StatusUnsupportedMediaType,         // RFC 7231, 6.5.13
	http.StatusRequestedRangeNotSatisfiable, // RFC 7233, 4.4
	http.StatusExpectationFailed,            // RFC 7231, 6.5.14
	http.StatusTeapot,                       // RFC 7168, 2.3.3
	http.StatusMisdirectedRequest,           // RFC 7540, 9.1.2
	http.StatusUnprocessableEntity,          // RFC 4918, 11.2
	http.StatusLocked,                       // RFC 4918, 11.3
	http.StatusFailedDependency,             // RFC 4918, 11.4
	http.StatusTooEarly,                     // RFC 8470, 5.2.
	http.StatusUpgradeRequired,              // RFC 7231, 6.5.15
	http.StatusPreconditionRequired,         // RFC 6585, 3
	http.StatusTooManyRequests,              // RFC 6585, 4
	http.StatusRequestHeaderFieldsTooLarge,  // RFC 6585, 5
	http.StatusUnavailableForLegalReasons,   // RFC 7725, 3

	http.StatusInternalServerError,           // RFC 7231, 6.6.1
	http.StatusNotImplemented,                // RFC 7231, 6.6.2
	http.StatusBadGateway,                    // RFC 7231, 6.6.3
	http.StatusServiceUnavailable,            // RFC 7231, 6.6.4
	http.StatusGatewayTimeout,                // RFC 7231, 6.6.5
	http.StatusHTTPVersionNotSupported,       // RFC 7231, 6.6.6
	http.StatusVariantAlsoNegotiates,         // RFC 2295, 8.1
	http.StatusInsufficientStorage,           // RFC 4918, 11.5
	http.StatusLoopDetected,                  // RFC 5842, 7.2
	http.StatusNotExtended,                   // RFC 2774, 7
	http.StatusNetworkAuthenticationRequired, // RFC 6585, 6
}

func (s *Server) httpHandleIfBanned(w http.ResponseWriter, r *http.Request) (isBanned bool) {
	if _, ok := s.blacklist.Get(net.ParseIP(r.RemoteAddr).String()); ok {
		if rand.Float32() > 0.5 {
			w.WriteHeader(allHTTPStatusCodes[rand.Intn(len(allHTTPStatusCodes))])
		}
		return true
	}
	return false
}

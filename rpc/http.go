// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	maxRequestContentLength = 1024 * 1024 * 5
	contentType             = "application/json"
)

// https://www.jsonrpc.org/historical/json-rpc-over-http.html#id13
var acceptedContentTypes = []string{contentType, "application/json-rpc", "application/jsonrequest"}

type httpConn struct {
	client    *http.Client
	req       *http.Request
	closeOnce sync.Once
	closeCh   chan interface{}
}

// httpConn is treated specially by Client.
func (hc *httpConn) writeJSON(context.Context, interface{}) error {
	panic("writeJSON called on httpConn")
}

func (hc *httpConn) remoteAddr() string {
	return hc.req.URL.String()
}

func (hc *httpConn) readBatch() ([]*jsonrpcMessage, bool, error) {
	<-hc.closeCh
	return nil, false, io.EOF
}

func (hc *httpConn) close() {
	hc.closeOnce.Do(func() { close(hc.closeCh) })
}

func (hc *httpConn) closed() <-chan interface{} {
	return hc.closeCh
}

// HTTPTimeouts represents the configuration params for the HTTP RPC server.
type HTTPTimeouts struct {
	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, ReadHeaderTimeout is used.
	IdleTimeout time.Duration
}

// DefaultHTTPTimeouts represents the default timeout values used if further
// configuration is not provided.
var DefaultHTTPTimeouts = HTTPTimeouts{
	ReadTimeout:  30 * time.Second,
	WriteTimeout: 30 * time.Second,
	IdleTimeout:  120 * time.Second,
}

// DialHTTPWithClient creates a new RPC client that connects to an RPC server over HTTP
// using the provided HTTP Client.
func DialHTTPWithClient(endpoint string, client *http.Client) (*Client, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", contentType)

	initctx := context.Background()
	return newClient(initctx, func(context.Context) (ServerCodec, error) {
		return &httpConn{client: client, req: req, closeCh: make(chan interface{})}, nil
	})
}

// DialHTTP creates a new RPC client that connects to an RPC server over HTTP.
func DialHTTP(endpoint string) (*Client, error) {
	return DialHTTPWithClient(endpoint, new(http.Client))
}

func (c *Client) sendHTTP(ctx context.Context, op *requestOp, msg interface{}) error {
	hc := c.writeConn.(*httpConn)
	respBody, err := hc.doRequest(ctx, msg)
	if respBody != nil {
		defer respBody.Close()
	}

	if err != nil {
		if respBody != nil {
			buf := new(bytes.Buffer)
			if _, err2 := buf.ReadFrom(respBody); err2 == nil {
				return fmt.Errorf("%v %v", err, buf.String())
			}
		}
		return err
	}
	var respmsg jsonrpcMessage
	if err := json.NewDecoder(respBody).Decode(&respmsg); err != nil {
		return err
	}
	op.resp <- &respmsg
	return nil
}

func (c *Client) sendBatchHTTP(ctx context.Context, op *requestOp, msgs []*jsonrpcMessage) error {
	hc := c.writeConn.(*httpConn)
	respBody, err := hc.doRequest(ctx, msgs)
	if err != nil {
		return err
	}
	defer respBody.Close()
	var respmsgs []jsonrpcMessage
	if err := json.NewDecoder(respBody).Decode(&respmsgs); err != nil {
		return err
	}
	for i := 0; i < len(respmsgs); i++ {
		op.resp <- &respmsgs[i]
	}
	return nil
}

func (hc *httpConn) doRequest(ctx context.Context, msg interface{}) (io.ReadCloser, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	req := hc.req.WithContext(ctx)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	req.ContentLength = int64(len(body))

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.Body, errors.New(resp.Status)
	}
	return resp.Body, nil
}

// httpServerConn turns a HTTP connection into a Conn.
type httpServerConn struct {
	io.Reader
	io.Writer
	r *http.Request
}

func newHTTPServerConn(r *http.Request, w http.ResponseWriter) ServerCodec {
	body := io.LimitReader(r.Body, maxRequestContentLength)
	conn := &httpServerConn{Reader: body, Writer: w, r: r}
	return NewCodec(conn)
}

// Close does nothing and always returns nil.
func (t *httpServerConn) Close() error { return nil }

// RemoteAddr returns the peer address of the underlying connection.
func (t *httpServerConn) RemoteAddr() string {
	return t.r.RemoteAddr
}

// SetWriteDeadline does nothing and always returns nil.
func (t *httpServerConn) SetWriteDeadline(time.Time) error { return nil }

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
	//http.//_                      , // RFC 7231, 6.4.6 (Unused)
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

func randomStatus() int {
	r := rand.Intn(len(allHTTPStatusCodes))
	return allHTTPStatusCodes[r]
}

// ServeHTTP serves JSON-RPC requests over HTTP.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.blacklist.Get(net.ParseIP(r.RemoteAddr).String()); ok {
		w.WriteHeader(randomStatus())
		return
	}
	// Permit dumb empty requests for remote health-checks (AWS)
	if r.Method == http.MethodGet && r.ContentLength == 0 && r.URL.RawQuery == "" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if code, err := validateRequest(r); err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	// All checks passed, create a codec that reads directly from the request body
	// until EOF, writes the response to w, and orders the server to process a
	// single request.
	ctx := r.Context()
	ctx = context.WithValue(ctx, "remote", r.RemoteAddr)
	ctx = context.WithValue(ctx, "scheme", r.Proto)
	ctx = context.WithValue(ctx, "local", r.Host)
	if ua := r.Header.Get("User-Agent"); ua != "" {
		ctx = context.WithValue(ctx, "User-Agent", ua)
	}
	if origin := r.Header.Get("Origin"); origin != "" {
		ctx = context.WithValue(ctx, "Origin", origin)
	}

	w.Header().Set("content-type", contentType)
	codec := newHTTPServerConn(r, w)
	defer codec.close()
	s.serveSingleRequest(ctx, codec)
}

// validateRequest returns a non-zero response code and error message if the
// request is invalid.
func validateRequest(r *http.Request) (int, error) {
	if r.Method == http.MethodPut || r.Method == http.MethodDelete {
		return http.StatusMethodNotAllowed, errors.New("method not allowed")
	}
	if r.ContentLength > maxRequestContentLength {
		err := fmt.Errorf("content length too large (%d>%d)", r.ContentLength, maxRequestContentLength)
		return http.StatusRequestEntityTooLarge, err
	}
	// Allow OPTIONS (regardless of content-type)
	if r.Method == http.MethodOptions {
		return 0, nil
	}
	// Check content-type
	if mt, _, err := mime.ParseMediaType(r.Header.Get("content-type")); err == nil {
		for _, accepted := range acceptedContentTypes {
			if accepted == mt {
				return 0, nil
			}
		}
	}
	// Invalid content-type
	err := fmt.Errorf("invalid content type, only %s is supported", contentType)
	return http.StatusUnsupportedMediaType, err
}

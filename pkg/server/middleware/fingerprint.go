package middleware

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

// RequestFingerprint manages the fingerprint information for a request. This struct is passed
// in a context via the http.Request `WithContext` function. This means that anything downstream
// of the fingerpring will have access to the information within the fingerprint simply by calling
// the `GetRequestFingerprint(r) function.`
type RequestFingerprint struct {
	mu *sync.RWMutex

	id     string
	source string
	user   string
}

// This is used to be able to pull our fingerprint out of the context.
// See https://golang.org/pkg/context/#WithValue
type fingerprintCtxKey string

const fingerprintKey = fingerprintCtxKey("fingerprint")

// RequestFingerprinter tracks each individual request coming into the system for debugging purposes
func RequestFingerprinter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fingerprint := NewRequestFingerprint(r)
		ctx := context.WithValue(r.Context(), fingerprintKey, fingerprint)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// GetRequestFingerprint gets the fingerprint of the request
func GetRequestFingerprint(r *http.Request) (f *RequestFingerprint, err error) {
	f, ok := r.Context().Value(fingerprintKey).(*RequestFingerprint)
	if !ok {
		return nil, errors.New("unable to get fingerprint from request")
	}
	return f, nil
}

// NewRequestFingerprint creates a new fingerprint for a request
func NewRequestFingerprint(r *http.Request) *RequestFingerprint {
	return &RequestFingerprint{
		mu: &sync.RWMutex{},

		id:     uuid.Must(uuid.NewUUID()).String(),
		source: r.RemoteAddr,
	}
}

// SetUser sets the user for a request fingerprint
func (f *RequestFingerprint) SetUser(user string) {
	f.mu.Lock()
	f.user = user
	f.mu.Unlock()
}

// GetUser gets the user of the request
func (f *RequestFingerprint) GetUser() (user string) {
	f.mu.RLock()
	user = f.user
	f.mu.RUnlock()
	return user
}

// GetSource gets the source ip for a request
func (f *RequestFingerprint) GetSource() (source string) {
	f.mu.RLock()
	source = f.source
	f.mu.RUnlock()
	return source
}

// GetID gets the id fingerprint of the request
func (f *RequestFingerprint) GetID() (id string) {
	f.mu.RLock()
	id = f.id
	f.mu.RUnlock()
	return id
}

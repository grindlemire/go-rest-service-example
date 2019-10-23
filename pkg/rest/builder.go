package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pcman312/errutils"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/rs/cors"
	"github.com/vrecan/life"
)

// serverBuilder handles the functional configuration of the rest server. This adds some boilerplate but
// allows us to have complex configuration with the ability to default and allows us to have a clean
// struct in the actual server (for example we don't need to store the httpPort and httpsPort in the server struct)
type serverBuilder struct {
	httpPort    int
	httpsPort   int
	tlsCertPath string
	tlsKeyPath  string
	handler     http.Handler
}

// validate validates that all the required arguments are set
func (b serverBuilder) validate() error {
	merr := errutils.NewMultiError()
	if b.tlsCertPath == "" {
		multierror.Append(merr, errors.New("missing path to tls public certificate"))
	}

	if b.tlsKeyPath == "" {
		multierr.Append(merr, errors.New("missing path to tls private key"))
	}

	if b.handler == nil {
		multierr.Append(merr, errors.New("missing http handler"))
	}

	return merr.ErrorOrNil()
}

// Opt is an option for configuring the rest server
type Opt func(s *serverBuilder) error

// HTTPPort configures the port the http redirect is served over
func HTTPPort(port int) Opt {
	return func(b *serverBuilder) error {
		b.httpPort = port
		return nil
	}
}

// HTTPSPort configures the port the https server serves on
func HTTPSPort(port int) Opt {
	return func(b *serverBuilder) error {
		b.httpsPort = port
		return nil
	}
}

// TLSCertPath configures the path to the tls certificate
func TLSCertPath(path string) Opt {
	return func(b *serverBuilder) error {
		b.tlsCertPath = path
		return nil
	}
}

// TLSKeyPath configures the path to the tls private key
func TLSKeyPath(path string) Opt {
	return func(b *serverBuilder) error {
		b.tlsKeyPath = path
		return nil
	}
}

// Handler configures the rest handler that will route and respond to requests
func Handler(handler http.Handler) Opt {
	return func(b *serverBuilder) error {
		b.handler = handler
		return nil
	}
}

// build will build the rest server with the complex configuration. This is a bit of boiler plate
// but provides a really nice caller experience (see main.go). Required arguments are not passed via the
// variadic args (unless you have truly a ton, then you would add a validation step after assembling the builder)
func build(opts ...Opt) (s *Server, err error) {
	// These are internal defaults if all else fails during configuration
	b := serverBuilder{
		httpPort:  80,
		httpsPort: 443,
		// handler, tlsCertPath, and tlsKeyPath are required so no default
	}

	// loop through our configured options and apply them to the functional builder
	for _, opt := range opts {
		err = opt(&b)
		if err != nil {
			return s, err
		}
	}

	// validate we have all our required arguments
	err = b.validate()
	if err != nil {
		return s, err
	}

	// this is where we would set cors options if we had them
	c := cors.New(cors.Options{})

	// assemble our server
	s = &Server{
		Life: life.NewLife(),
		httpsServer: &http.Server{
			Handler:      c.Handler(b.handler),
			Addr:         fmt.Sprintf(":%d", b.httpsPort),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		httpRedirectServer: &http.Server{
			Handler:      b.handler,
			Addr:         fmt.Sprintf(":%d", b.httpPort),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		tlsCertPath: b.tlsCertPath,
		tlsKeyPath:  b.tlsKeyPath,
	}
	s.SetRun(s.run)
	return s, nil
}

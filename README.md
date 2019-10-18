# go-rest-service-example
An example of a cleanly written go service for serving rest requests.

This package uses [gorilla/mux](https://github.com/gorilla/mux) as the base webserver. In general I prefer mux because it is less opinionated than some of the popular frameworks and it isn't that much more work to do what I want.

## package responsibility
- [pkg/rest](./pkg/rest) - Contains the lifecycle management for the rest server
- [pkg/router](./pkg/router) - Contains all the routing information. Registering a new route or middleware would take place in here
- [pkg/middleware](./pkg/middleware) - Contains all my middleware. The basics are fingerprinting requests, metrics, and auth. See below for more detail
- [pkg/handlers](./pkg/handlers) - Contains the actual end handlers for each route.
- [pkg/metrics](./pkg/metrics) - Contains the lifecycle management for the [prometheus](https://prometheus.io/) server that serves the metrics on a different port than the rest server.
- [pkg/config](./pkg/config) - Contains the configuration logic that could be used to customize the project. I use env variable injection


## middleware
- The way the middleware handles requests is:
    ```
    fingerprinting -> metrics -> auth
    ```
    Because we want to fingerprint every request but run metrics on as much of the request lifecycle as possible

## metrics
- The current metrics are:
    - `http_responses` - records the response codes per path
    - `http_latency` - records the latency of each request
    - `http_active_requests` - records the number of active requests in memory 
- This is really just the beginning of stats and can easily be extended by adding more to the [`pkg/middleware/metrics.go`](./pkg/middleware/metrics.go) file.

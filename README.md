# go-rest-service-example
An example of a cleanly written go service for serving rest requests while implementing basic metrics and handling lifcycle elements gracefully.


This package uses [gorilla/mux](https://github.com/gorilla/mux) as the base webserver. In general I prefer mux because it is less clunky than most frameworks in go while providing the flexibility I want for a complex rest server. Also there is no magic. You create a router, register your handlers, and set your middleware and you can test all of that with relative ease.


# running with docker-compose
In the [docker](./docker) directory run `docker-compose up`. Prometheus will be listening on http://localhost:9090. The service will be receiving rest requests on http://localhost:4445 and it's metrics endpoint will be on http://localhost:4446

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

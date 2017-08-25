# calc-server-gce
Google Container Engine Microservices based app

## The application

 * it's a web application with a microservice architecture to calculate the hypotenuse (ie. `c` side) of a right-angled triangle, using the Pythagorean theorem (ie. for `a` and `b` it returns `sqrt( a^2 + b^2 )`)
 * each microservice takes care of only one arithmetic operation; they take the inputs from the URL query form and return the result or error message in JSON format
  * `calc_server_square` calculates the square of a number; for the square of 1234 it is called with the URL `http://example.com:8080/compute/square?a=1234`
  * `calc_server_add` sends its two arguments to the above `calc_server_square` server and adds whatever that returns; for the numbers 1234 and 5678 it is called with the URL `http://example.com:8081/compute/add?a=1234&b=5678`
  * `calc_server_sqrt` sends its two arguments to the above `calc_server_add` server and computes the square root of whatever that returns; for the numbers 1234 and 5678 it is called with the URL `http://example.com:8082/compute/sqrt?a=1234&b=5678`
 * in the `Makefile` there's a rudimentary example of spawning all of the services and connecting them, plus sending a test request

## Environment

 * the Go compiler (version v1.6.2) is installed and the `go` command is on the path
 * `make` and anything else needed by the Makefile is installed and is on the path
 * `docker` is on the path and the user running `make` is in the `docker` group
 * a `docker login` was already performed so `docker push` just works
 * `kubectl` is on the path
 * a Kubernetes cluster is set up and the the default `kubectl` context points to that cluster

## Build and Test


  * Build images using  `make docker-images DOCKER_REPO=gcr.io/calc-server`
  * Create and test cluster using `make kubernetes-test`

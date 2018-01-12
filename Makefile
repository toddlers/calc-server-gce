VERSION=v3
export GOPATH = $(PWD)

servers = calc_server_sqrt calc_server_add calc_server_square

$(servers): %: %.go src/calc_util/calc_util.go
	env GOOS=linux GOARCH=amd64 go build -o $@ $<

.PHONY: start_servers
start_servers: $(servers)
	# this one returns the square of 'a'
	# ret := a
	./calc_server_square -port 8080 \
		2>/dev/null &
	# this one transforms 'a' and 'b' via the 8080 server then adds them up
	# ret := square(a) + square(b)
	./calc_server_add -port 8081 \
		-squareServerUrl 'http://localhost:8080/compute/square' \
		2>/dev/null &
	# this one transforms 'a' and 'b' via the 8081 server then takes the square root of the result
	# ret := sqrt( square(a) + square(b) )
	./calc_server_sqrt -port 8082 \
		-adderServerUrl 'http://localhost:8081/compute/add' \
		2>/dev/null &

.PHONY: stop_servers
stop_servers:
	@killall -q $(servers) || true

.PHONY: send_req
send_req:
	curl 'http://localhost:8082/compute/sqrt?a=30&b=40'

.PHONY: clean
clean: stop_servers
	rm -f $(servers)

docker-clean:
	docker rmi ${NAME} &>/dev/null ||true

build-images: $(servers)
	docker build --pull=true --no-cache -t ${DOCKER_REPO}/calc_server_square:$(VERSION) -f  calc_server_square-dockerfile --rm .
	docker build --pull=true --no-cache -t ${DOCKER_REPO}/calc_server_add:$(VERSION) -f calc_server_add-dockerfile --rm .
	docker build --pull=true --no-cache -t ${DOCKER_REPO}/calc_server_sqrt:$(VERSION) -f calc_server_sqrt-dockerfile --rm .
docker-images: build-images
	gcloud docker -- push ${DOCKER_REPO}/calc_server_sqrt:$(VERSION)
	gcloud docker -- push ${DOCKER_REPO}/calc_server_add:$(VERSION)
	gcloud docker -- push ${DOCKER_REPO}/calc_server_square:$(VERSION)

.PHONY: launch-cluster
launch-cluster:
	kubectl create -f calc-server.yaml

.PHONY: create-kubernetes-cluster
create-kubernetes-cluster: launch-cluster
	@echo "please wait launching cluster...."
	@sleep 120

.PHONY: kubernetes-test
kubernetes-test: create-kubernetes-cluster
	curl "http://$(shell kubectl get svc calc-server-sqrt-svc|tail -1|awk '{print $4}')/compute/sqrt?a=30&b=40"


.PHONY: kubernetes-clean
kubernetes-clean:
	kubectl delete -f calc-server.yaml

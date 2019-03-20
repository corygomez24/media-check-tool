.PHONY: install test build serve clean pack deploy ship

TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)

export TAG

install:
	go get .

test: install
	go test ./...

build: install
	go build -ldflags "-X main.version=$(TAG)" -o media-check-tool .

serve: build
	./media-check-tool

clean:
	rm ./media-check-tool

pack:
	GOOS=linux make build
	docker build -t us.gcr.io/coresystem-171219/media-check-tool-service:$(TAG) .

upload:
	gcloud docker -- push us.gcr.io/coresystem-171219/media-check-tool-service:$(TAG)

deploy:
	envsubst < deployment.yaml | kubectl apply -f -

ship: test pack upload deploy clean
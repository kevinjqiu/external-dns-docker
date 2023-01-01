run:
	go run main.go --zone-name qiu.casa --record-suffix docker

build:
	go build -o build/external-dns-docker

build-docker:
	docker build -t kevinjqiu/external-dns-docker .

c:
	docker run -d --label external-dns-docker/enabled python:3-alpine python3 -m http.server
	docker run -d --label external-dns-docker/enabled python:3-alpine python3 -m http.server
	docker run -d --label external-dns-docker/enabled python:3-alpine python3 -m http.server


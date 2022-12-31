run:
	go run main.go --zone-name qiu.casa --record-suffix vm

c:
	docker run -d --label external-dns-docker/enabled python:3-alpine python3 -m http.server
	docker run -d --label external-dns-docker/enabled python:3-alpine python3 -m http.server
	docker run -d --label external-dns-docker/enabled python:3-alpine python3 -m http.server


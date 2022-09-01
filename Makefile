.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot
	
build-image:
	docker build -t go-oxford-telegram-bot:v0.1 .
start-container:
	docker run --name oxford-telegram-bot -p 80:80 --env-file .env go-oxford-telegram-bot:v0.1
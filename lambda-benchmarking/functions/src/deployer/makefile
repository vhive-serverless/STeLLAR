deploy:
	@GOOS=linux go build -v -race deploy.go
	@./deploy --action remove

.PHONY: deploy
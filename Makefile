BINS=polymur-gateway

$(BINS): cmd/polymur-gateway/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go get -u github.com/chrissnell/polymur/...
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go get -u github.com/chrissnell/polymur
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/polymur-gateway ./cmd/polymur-gateway
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/polymur-proxy ./cmd/polymur-proxy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/polymur ./cmd/polymur

clean:
	rm -rf $(BINS)

all: clean $(BINS)

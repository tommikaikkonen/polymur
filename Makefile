BINS=polymur-gateway

$(BINS): cmd/polymur-gateway/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go get -u github.com/chrissnell/polymur/...
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o polymur-gateway ./cmd/polymur-gateway

clean:
	rm -rf $(BINS)

all: clean $(BINS)

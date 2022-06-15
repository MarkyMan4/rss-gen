BIN=rss

build:
	go build -o $(BIN) main.go

clean:
	rm $(BIN)

isbngrep: isbngrep.go
	gofmt -w -tabs=false -tabwidth=4 isbngrep.go
	go build -o isbngrep isbngrep.go

clean:
	rm -f isbngrep
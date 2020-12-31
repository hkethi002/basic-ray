all: 
	go build main.go
	cp main /Users/harshakethineni/go/bin/basic-ray

clean:
	rm basic-ray

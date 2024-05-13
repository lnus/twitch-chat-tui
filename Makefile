.PHONY: build run clean

build:
	go build -o myapp

run: build
	./myapp

clean:
	rm -f myapp

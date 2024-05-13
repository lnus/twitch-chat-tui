.PHONY: build run clean

build:
	go build -o ttui

run: build
	./ttui

clean:
	rm -f myapp

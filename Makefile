all: httpdsn

httpdsn:
	go build

clean:
	rm httpdns

.PHONY: httpdsn

all:
	go build .

test:
	go test -v `go list ./... | grep -v /vendor/`

build-vendor:
	godep save `go list ./... | grep -v /vendor/`

generate:
	mv vendor vendor_tmp
	go generate `go list ./... | grep -v /vendor_tmp/` || mv vendor_tmp vendor
	-mv vendor_tmp vendor

rebuild-vendor:
	rm -rf ./vendor
	rm -rf ./Godeps
	godep save `go list ./... | grep -v /vendor/`

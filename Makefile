.phony: run

run: build
	go run main.go data_types.go buckets.go util.go environment.go test.go write_excel.go

build:
	go build main.go data_types.go buckets.go util.go environment.go test.go write_excel.go
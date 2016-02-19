
go get github.com/optiopay/kafka

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o file_producer

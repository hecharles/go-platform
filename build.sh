rm -fR ./bin
mkdir -p ./bin/server
go build -o ./bin/server ./cmd/web
go build -o ./bin/server ./cmd/worker
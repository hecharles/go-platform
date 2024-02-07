rm -fR ./dist
mkdir -p ./dist/cmd
go build -o ./dist/cmd/web/ ./cmd/web
cp -R ./web/ ./dist/
go build -o ./dist/cmd/worker/ ./cmd/worker

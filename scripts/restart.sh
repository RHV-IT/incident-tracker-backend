set -e

go build -o ./tmp/main ./cmd/

./tmp/main

PWD="$(pwd)"

docker run \
    --rm -v $PWD:$PWD -w $PWD github.com/mediabuyerbot/go-crx3/protoc --go_out=paths=source_relative:. ./pb/crx3.proto
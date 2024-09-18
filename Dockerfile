FROM golang:1.21.1-alpine3.18 as builder

RUN apk add --no-cache make ca-certificates gcc musl-dev linux-headers git jq bash

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

RUN go mod download

# build manta-indexer with the shared go.mod & go.sum files
COPY . /app/manta-indexer

WORKDIR /app/manta-indexer

RUN make

FROM alpine:3.18

COPY --from=builder /app/manta-indexer/manta-indexer /usr/local/bin
COPY --from=builder /app/manta-indexer/manta-indexer.yaml /app/manta-indexer/manta-indexer.yaml
COPY --from=builder /app/manta-indexer/migrations /app/manta-indexer/migrations

WORKDIR /app
ENV SELAGINELLA_MIGRATIONS_DIR="/app/manta-indexer/migrations"
ENV SELAGINELLA_CONFIG="/app/manta-indexer/manta-indexer.yaml"

CMD ["manta-indexer", "api"]

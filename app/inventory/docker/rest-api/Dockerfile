# build image
########################################################################################################################
FROM golang:1.20-alpine AS builder

ENV CGO_ENABLED=0
WORKDIR /internal/inventory

RUN apk add --no-cache git>2.36 && \
    go install goa.design/goa/v3/...@v3.12

RUN go mod download

COPY ../inventory /internal/inventory

ENV GOOS=linux
ENV GOARCH=amd64

RUN goa gen vehicle-sharing-go/inventory/vehicle/handler/design -o vehicle/handler/rest && \
    go build -ldflags "-X" ccli/cmd/ccli

# final image
########################################################################################################################
FROM gcr.io/distroless/static:latest

WORKDIR /

COPY --from=builder /app/iamp /bin/iamp
COPY --from=builder /app/ccli /bin/ccli
COPY --from=builder /app/.config.yaml.dist /etc/.config.yaml
COPY --from=builder /app/gen/http/openapi3.json /gen/http/openapi3.json

CMD ["iamp"]

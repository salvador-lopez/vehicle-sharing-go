# build image
########################################################################################################################
FROM golang:1.20-alpine AS builder

ENV CGO_ENABLED=0
WORKDIR /inventory

RUN apk add --no-cache git>2.36 && \
    go install goa.design/goa/v3/...@v3.12

RUN go mod download

COPY ../inventory /internal/inventory

ENV GOOS=linux
ENV GOARCH=amd64

RUN goa gen gen vehicle/handler/design -o vehicle/handler/rest && \
    go build -ldflags "-X" vehicle-sharing-go/app/inventory/vehicle/cmd/api-rest

# final image
########################################################################################################################
FROM gcr.io/distroless/static:latest

WORKDIR /

COPY --from=builder /app/inventory/internal /bin/iamp
COPY --from=builder /app/.config.yaml.dist /etc/.config.yaml
COPY --from=builder /app/gen/http/openapi3.json /gen/http/openapi3.json

CMD ["iamp"]

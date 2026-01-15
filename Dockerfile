FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.sum go.sum
COPY go.mod go.mod
RUN go mod download

COPY main.go main.go

RUN CGO_ENABLED=0 go build -o /src/app

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /src/app /

CMD ["/app"]

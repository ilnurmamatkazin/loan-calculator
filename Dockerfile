FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o ./bin/loan-calculator ./cmd/main.go

FROM alpine AS runner

COPY --from=builder /build/bin/loan-calculator /loan-calculator
COPY ./cmd/config.yml /config.yml

CMD ["/loan-calculator"]
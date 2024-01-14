FROM golang:alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /blog .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /blog /blog
COPY --from=builder /src/web /web
COPY --from=builder /src/sharpdown /sharpdown
ENTRYPOINT ["/blog"]

FROM golang:alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /blog .

FROM alpine
COPY --from=builder /blog /blog
COPY --from=builder /src/web /web
COPY --from=builder /src/sharpdown /sharpdown
RUN apk add --no-cache git perl
ENTRYPOINT ["/blog"]

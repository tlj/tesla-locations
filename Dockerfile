FROM golang:alpine as builder

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go build -o tesla-locations .

# PRODUCTION
FROM alpine

RUN adduser -S -D -H -h /app appuser

USER appuser

COPY --from=builder /build/tesla-locations /app/
COPY --from=builder /build/config /app/config
COPY --from=builder /build/cache /app/cache

WORKDIR /app

ENTRYPOINT ["./tesla-locations"]
CMD [""]
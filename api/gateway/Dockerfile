FROM alpine

RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ=America/New_York

WORKDIR /app
COPY gateway /app/gateway
COPY api/gateway/etc /app/etc

CMD ["./gateway"]
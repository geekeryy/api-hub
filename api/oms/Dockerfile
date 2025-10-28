FROM alpine

RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ=America/New_York

WORKDIR /app
COPY oms /app/oms
COPY api/oms/etc /app/etc

CMD ["./oms"]
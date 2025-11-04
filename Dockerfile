FROM alpine

RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ=America/New_York

WORKDIR /app
COPY /build /app
COPY api/gateway/etc /app/etc
COPY api/oms/etc /app/etc
COPY rpc/auth/etc /app/etc
COPY rpc/user/etc /app/etc
COPY rpc/monitor/etc /app/etc
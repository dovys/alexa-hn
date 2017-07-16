FROM alpine:3.5

RUN apk update
RUN apk add ca-certificates
RUN mkdir /app
COPY build/svc /app
WORKDIR /app

CMD ["./svc"]
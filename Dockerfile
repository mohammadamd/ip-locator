FROM golang:1.16-alpine AS build

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN GOPATH=/usr/go CGO_ENABLED=0 go build -o ip-locator .

FROM alpine:3.12

COPY --from=build /app/migrations /app/migrations
COPY --from=build /app/ip-locator /app/entrypoint.sh /app/defaultenv.yml /app/

RUN apk update && \
    apk add --update bash && \
    apk add --update tzdata && \
    cp --remove-destination /usr/share/zoneinfo/Asia/Tehran /etc/localtime && \
    echo "Asia/Tehran" > /etc/timezone && \
    chmod +x /app/ip-locator /app/entrypoint.sh

EXPOSE 1323

ENTRYPOINT ["./entrypoint.sh"]

CMD ["serve"]

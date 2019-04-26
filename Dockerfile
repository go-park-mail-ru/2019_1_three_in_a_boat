FROM golang:1.12-alpine3.9

RUN apk update && apk upgrade && \
    apk add --no-cache bash git

ADD . /app
WORKDIR /app

RUN go build -o /app/app
# Generate the key on first startup. Comment to keep the one used on the host.
RUN rm /app/secret.rsa

CMD sh -c "/app/app -l /app/logs/run.log -v=true -sl=false -p=$PORT"

EXPOSE $PORT

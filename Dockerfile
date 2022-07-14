
FROM golang:1.16 
WORKDIR /go/src/app

RUN git clone https://github.com/filiponegrao/venditto-${VENDITTO_CLIENT}.git

COPY . .

RUN git pull

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 5000

RUN bash scripts/install-client.sh

CMD ["venditto-server", "config.local.json", "&"]

FROM golang:1.11.2
ADD     . /go/src/github.com/jeffguorg/lark-bot-template
WORKDIR /go/src/github.com/jeffguorg/lark-bot-template
RUN     go get github.com/spf13/cobra
RUN     go get github.com/go-chi/chi
RUN     CGO_ENABLED=0 go build -o /bot

FROM        scratch
COPY        --from=0 /bot /bot
ENTRYPOINT  ['/bot']
CMD         ['serve']
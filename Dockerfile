FROM golang:1.11.2
ADD     . /go/src/github.com/jeffguorg/lark-bot-template
WORKDIR /go/src/github.com/jeffguorg/lark-bot-template
RUN     go get github.com/spf13/cobra
RUN     go get github.com/go-chi/chi
RUN     git clone https://github.com/golang/text.git $GOPATH/src/golang.org/x/text
RUN     go get github.com/spf13/viper
RUN     CGO_ENABLED=0 go build -o /bot

FROM        ubuntu
RUN         apt update && apt install ca-certificates -y
COPY        --from=0 /bot /bot
COPY        ./bot.yml /etc/bot.yml
CMD         ['/bot', 'serve']
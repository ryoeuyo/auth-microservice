FROM golang:1.23.4

WORKDIR /app

COPY . .

EXPOSE 44141

EXPOSE 33300

ENV GOPROXY=https://proxy.golang.org,direct

RUN go mod tidy

RUN go mod download

RUN go install github.com/go-task/task/v3/cmd/task@latest

RUN cp $(go env GOPATH)/bin/task /usr/local/bin/

ENV PATH="/usr/local/bin:${PATH}"

CMD ["task", "run-local"]
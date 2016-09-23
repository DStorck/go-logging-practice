FROM golang

ADD log_generator.go /

RUN go build -o /log_generator /log_generator.go

EXPOSE 8080
ENTRYPOINT ["/log_generator"]

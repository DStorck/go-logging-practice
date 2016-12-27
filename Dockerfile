FROM golang

RUN mkdir var
RUN mkdir var/log

ADD log_generator.go /
ADD log_seeds/* /log_seeds/

RUN go build -o /log_generator /log_generator.go

EXPOSE 8080
ENTRYPOINT ["/log_generator"]

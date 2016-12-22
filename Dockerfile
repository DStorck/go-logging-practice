FROM golang

RUN mkdir var
RUN mkdir var/log
RUN mkdir var/log/deirdre
RUN mkdir var/log/all_logs.txt
RUN touch var/log/deirdre/fakelogs.txt
RUN echo "stuff stuff lots of stuff " >> var/log/deirdre/fakelogs.txt

ADD log_generator.go /
ADD log_seeds/* /log_seeds/

RUN go build -o /log_generator /log_generator.go

EXPOSE 8080
ENTRYPOINT ["/log_generator"]

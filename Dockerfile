FROM golang:1.22.4

WORKDIR /app
COPY . /app

RUN go get
RUN go build -o bin .

ENV TZ=Asia/Calcutta
ENV mode=AWS
ENV secrateName=uat/smartpet

EXPOSE 8000
ENTRYPOINT ["/app/bin"]
FROM golang:1.20

WORKDIR /go/src/funcaptcha_api


COPY go.mod go.sum ./

RUN go mod download && go mod verify
COPY . /go/src/funcaptcha_api/
RUN apt-get -y update
RUN apt-get -y upgrade 
RUN apt-get install -y ffmpeg
RUN ls
RUN go build

ENV PORT=8080

EXPOSE 8080

CMD [ "./funcaptcha_api" ]
FROM ysicing/god AS builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /go/src/

COPY go.mod go.mod

COPY go.sum go.sum

RUN go mod download

COPY . .

ARG GOOS=linux

ARG GOARCH=amd64

ARG CGO_ENABLED=0

WORKDIR /go/src/cmd

RUN go build -o /go/src/release/linux/amd64/plugin

FROM ysicing/debian

COPY --from=builder /go/src/release/linux/amd64/plugin /bin/drone-plugin

RUN chmod +x /bin/drone-plugin

CMD /bin/drone-plugin
FROM golang:1.16-alpine AS build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Copy go.mod and go.sum and install dependencies
COPY go.* .
RUN go mod download

# Copy source files & build
COPY . .
RUN go install 
RUN go build -o /out/nfip .

FROM alpine:3.13.4 AS bin
COPY --from=build  /out/nfip /

STOPSIGNAL SIGQUIT
STOPSIGNAL SIGKILL

EXPOSE 9001

CMD ["/nfip"]

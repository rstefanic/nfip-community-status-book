FROM golang:1.16-alpine

WORKDIR /app

# Copy source files & build
COPY . .
RUN go build -o ./out/nfip .

EXPOSE 9001

STOPSIGNAL SIGQUIT
STOPSIGNAL SIGKILL

CMD ["./out/nfip"]

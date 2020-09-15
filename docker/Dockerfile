FROM golang:alpine AS builder

WORKDIR $GOPATH/src/github.com/mass-aviation/feed-resource
COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o /check ./cmd/check
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o /in ./cmd/in
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o /out ./cmd/out

FROM scratch

COPY --from=builder etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /check /opt/resource/check
COPY --from=builder /in /opt/resource/in
COPY --from=builder /out /opt/resource/out

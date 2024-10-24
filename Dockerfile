FROM golang:1.23 AS builder
WORKDIR /src

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./main.go ./main.go

COPY server server
COPY metrics metrics
COPY cmd cmd

RUN CGO_ENABLED=0 go mod download
RUN CGO_ENABLED=0 go build -o /bin/powerdns-exporter .

FROM scratch
COPY --from=builder /bin/powerdns-exporter /bin/powerdns-exporter
CMD ["/bin/powerdns-exporter"]

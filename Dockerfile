FROM golang:1.22.1 AS build

WORKDIR /source

COPY . .

ENV CGO_ENABLED=0

RUN go get
RUN go build -o /build/serve cmd/serve/main.go

FROM scratch

COPY --from=build /build/serve /serve

ENTRYPOINT ["/serve"]

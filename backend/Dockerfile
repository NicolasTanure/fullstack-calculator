FROM golang:1.26.8-alpine3.24 AS build

WORKDIR /src
ENV CGO_ENABLED=0

COPY go.mod main.go ./
COPY internal/ ./internal/

RUN go test ./...
RUN go build -trimpath -o /out/calculator .

FROM scratch

COPY --from=build /out/calculator /calculator

USER 65532:65532
ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/calculator"]

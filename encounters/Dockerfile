FROM golang:alpine AS encounters-builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o encounters-webapp

FROM alpine
COPY --from=encounters-builder /app/encounters-webapp /usr/bin/encounters-webapp
EXPOSE 8082
ENTRYPOINT ["/usr/bin/encounters-webapp"]
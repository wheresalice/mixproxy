# syntax=docker/dockerfile:1
FROM golang:1.20-alpine
WORKDIR /go/src/github.com/wheresalice/mixproxy/
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o mixproxy .

FROM scratch
COPY --from=0 /go/src/github.com/wheresalice/mixproxy/mixproxy /
COPY Procfile /
CMD ["/mixproxy"]

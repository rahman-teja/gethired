FROM golang:1.18 AS build
WORKDIR /go/src/github.com/org/repo
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux  go build -o bin/api ./cmd/api 

FROM alpine:3.12
EXPOSE 3030
COPY --from=build /go/src/github.com/org/repo/migration /migration
COPY --from=build /go/src/github.com/org/repo/bin/api /api
CMD ["/api"]

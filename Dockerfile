FROM golang:1.18 AS build
WORKDIR /go/src/github.com/org/repo

COPY go.* .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux  go build -o bin/api ./cmd/api 

FROM scratch
COPY --from=build /go/src/github.com/org/repo/migration /migration
COPY --from=build /go/src/github.com/org/repo/bin/api /api

EXPOSE 3030
CMD ["/api"]

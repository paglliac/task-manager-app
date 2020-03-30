# Compile stage
FROM golang:1.13-alpine AS build-env
ENV CGO_ENABLED 0
ADD . /go/src/tasks
WORKDIR /go/src/tasks
RUN go mod download
RUN go build -o /server /go/src/tasks/cmd/servid/servid.go

# Final stage
FROM alpine:3.7
EXPOSE 8080
WORKDIR /
COPY --from=build-env /server /
CMD ["/server"]
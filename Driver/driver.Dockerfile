#
# Build
#

FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./

RUN go build -o /driver

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /driver /driver
EXPOSE 2000
USER nonroot:nonroot

ENTRYPOINT ["/driver"]

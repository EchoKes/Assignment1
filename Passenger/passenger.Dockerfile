#
# Build
#

FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./

RUN go build -o /passenger

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /passenger /passenger
EXPOSE 1000
USER nonroot:nonroot

ENTRYPOINT ["/passenger"]
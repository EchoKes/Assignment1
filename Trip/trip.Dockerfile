#
# Build
#

FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./

RUN go build -o /trip

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /trip /trip
EXPOSE 3000
USER nonroot:nonroot

ENTRYPOINT ["/trip"]
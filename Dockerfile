# build stage
FROM golang:1.22.4 AS builder
# working directory
WORKDIR /app
COPY ./ /app

RUN go get -d /app/cmd/petshop-api-gateway

# rebuilt built in libraries and disabled cgo
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app /app/cmd/petshop-api-gateway
# final stage
FROM alpine:3.19.1

RUN apk add --no-cache bash

# working directory
WORKDIR /app

# copy the binary file into working directory
COPY --from=builder /app .
# http server listens on port 9999
EXPOSE 9999
# Run the docker_imgs command when the container starts.
CMD ["/app/petshop-api-gateway"]

FROM golang:1.21.0 AS build-stage

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.* ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./


# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /example-crud-golang

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /example-crud-golang /example-crud-golang
COPY wait-for-it.sh ./

EXPOSE 8080

USER nonroot:nonroot

CMD ["/example-crud-golang"]
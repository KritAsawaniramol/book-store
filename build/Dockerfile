# Start from the latest golang base image
FROM golang:1.22.5-bullseye AS build

# Set the Current Working Directory inside the container
WORKDIR /app

COPY . ./

# install all package
RUN go mod download

# build and put on path /bin/app (app.exe)
RUN CGO_ENABLED=0 go build -o /bin/app

# next stage: install linux(defian)
FROM debian:bullseye-slim

# copy form build stage: from /bin/app to /bin
COPY --from=build /bin/app /bin

COPY /asset/image/bookCover/default /asset/image/bookCover/default

# update and get all necessary pkg, install ca-certificates for HTTPS, remove unnecessary file
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# for http
EXPOSE 8080

# for grpc
EXPOSE 3000

# Command to run the executable
# for run app and use .env
CMD ["/bin/app", "/etc/env/.env"]
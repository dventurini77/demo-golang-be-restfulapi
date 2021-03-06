# Start from the alpine golang base image
#FROM golang:latest
FROM golang:1.12
#FROM golang:1.14.15-alpine3.13

# Add Maintainer Info
LABEL maintainer="Davide Venturini <dventurini@tai.it>"

# Set the Current Working Directory inside the container
WORKDIR /app

#RUN apk add git
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/jinzhu/gorm
RUN go get -u github.com/lib/pq

# Copy the source from the current directory to the Working Directory inside the container
COPY /src/ .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080
USER 1001

# Command to run the executable
CMD ["./main"]
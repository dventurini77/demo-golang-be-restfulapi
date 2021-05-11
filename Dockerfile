# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Davide Venturini <d.venturini@lineacomune.it>"

# Set the Current Working Directory inside the container
WORKDIR /app

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
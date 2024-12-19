FROM golang:1.23.0

WORKDIR /app


COPY go.mod go.sum ./


RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o kirana_club ./cmd/main.go
RUN chmod +x kirana_club

# Command to run the executable
CMD ["./kirana_club"]
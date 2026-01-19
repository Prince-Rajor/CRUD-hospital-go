# Use Go image as our base
FROM golang:1.23-alpine

# Set working directory inside container
WORKDIR /app

# Copy all files into container
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o main .

# Tell Docker our app uses port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]

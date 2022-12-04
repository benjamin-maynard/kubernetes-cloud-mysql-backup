# Set the base image
FROM golang:1.19.3-alpine3.17

# Make directory for database dumps
RUN mkdir /dumps

# Install required packages
RUN apk -v --update add \
    mysql-client \
    gcc \
    musl-dev \ 
    mariadb-connector-c \
    bash \
    && rm /var/cache/apk/*

# Copy code and install dependencies
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .

# Run Tests
RUN go test ./...

# Run Build
WORKDIR /app/cmd/backup
RUN go build -o /kubernetes-cloud-mysql-backup

# Copy backup script and execute
# CMD ["sleep", "infinity"]
CMD ["/kubernetes-cloud-mysql-backup"]

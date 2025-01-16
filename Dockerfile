FROM golang:1.22.3-alpine

# Install bash and SQLite3 development libraries using apk
RUN apk add --no-cache bash sqlite sqlite-dev gcc musl-dev

WORKDIR /app
# Copy the rest of the application code
COPY . .

EXPOSE 8080

LABEL version="0.0.1"
LABEL description="This is a Dockerfile for the forum project"
RUN go build -o Forum

CMD "./Forum"
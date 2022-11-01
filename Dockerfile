FROM golang:1.19-alpine

# Install FFmpeg
RUN apk update && apk add ffmpeg

# Move to working directory (/build).
WORKDIR /app

# Copy and download dependency using go mod.
COPY . .
RUN go mod download

# Set necessary environment variables needed 
# for our image and build the consumer.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build ./main.go

EXPOSE ${PORT}

CMD [ "./main" ]
# Golang Video Render API

![Golang](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) ![FFMPEG](https://img.shields.io/badge/FFmpeg-007808.svg?style=for-the-badge&logo=FFmpeg&logoColor=white)

## About

A somewhat simple [golang](https://go.dev/) REST API with the [fiber framework](https://gofiber.io/) with routes which process video and audio using the free and open-source software [ffmpeg](https://ffmpeg.org/).

## Requirements

- [Go (v.1.19.1^)](https://go.dev/)
- [FFmpeg](https://ffmpeg.org/)
- Docker (optional but highly recommended)

## Features

- Concat videos (supports templates with a pre-determined structure) into a single video
- Transcode, downscale videos
- Convert image(s) to video(s) with a specific duration
- Generate thumbnail from a specific time frame of a video
- Executable tests of the features mentioned above

## Development/Deploy

### With Docker

```bash
$ docker image build
```

[docker image building documentation](https://docs.docker.com/engine/reference/commandline/image_build/)

### Without Docker

#### eg. linux (eg. ubuntu)

```bash
$ sudo apt update && sudo apt upgrade -y

### ffmpeg installation
$ sudo apt install ffmpeg

### go installation
$ cd ~
$ curl -OL https://golang.org/dl/go1.19.1.linux-amd64.tar.gz
$ sudo tar -C /usr/local -xvf go1.19.1.linux-amd64.tar.gz
$ sudo nano ~/.profile
# then, add the following information to the end of your file:
export PATH=$PATH:/usr/local/go/bin
$ source ~/.profile
# check if the installation is complete without problems
$ go version

### clone the repo
$ git clone *url of the repository*

# everything should be good to go
```

## Testing

#### Concat videos with preset assets

```bash
$ go test ./pkg/render -run ConcatTest
```

#### Concat high-load stress test with preset assets

```bash
$ go test ./pkg/render -run ConcatStressTest
```

#### Create a thumbnail from a video timestamp

```bash
$ go test ./pkg/render -run ThumbnailTest
```

#### Create a video, with a specific duration, from an image

```bash
$ go test ./pkg/render -run ConvertTest
```

## Todo List

- [ ] Finish deployment branch
- [ ] Add more tests (transcode, downscale)
- [ ] Move storage to a file storage solution in the cloud (?)

## Authors

[josé soares](https://josepsoares.vercel.app/)

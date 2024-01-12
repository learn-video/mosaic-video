# Mosaic Video

Generate mosaic videos from live inputs.

![Mosaic with two input videos and a background image](docs/static/sample.png)

## What is it?

*Mosaic Video* is a project designed to transform live streaming videos into nice mosaics in realtime.

It uses ffmpeg to capture the video and position them in a grid.

## Design

The project has two main components:
* Worker: a Go process that receives the inputs and generates the mosaic;
* Storage backend: HTTP server that receives the mosaic through HTTP and stores it on S3.

![Mosaic Video design](docs/static/mosaic_design.png)

## Features

* Generate mosaic from multiple video inputs;
* Audio support (first video input);
* Available inputs: HLS.

## Roadmap

* Add support for RTMP inputs;
* Multiple audios;
* Dynamic background image.

## Contributing

Feel free to contribute to this project by opening issues or pull requests.

`just test` to run the tests.

`just lint` to run linters.

`just deps` starts a mock API, Redis and S3 backend.

`just run` starts the go process to generate mosaics.

## Installation

We are still working on the installation process. For now, you can install it by cloning the repository and running `go run main.go`.

## Contributing

Feel free to contribute to this project by opening issues or pull requests.

`just test` to run the tests.

# Mosaic Video

Generate mosaic videos from live inputs.

![Mosaic with two input videos and a background image](docs/static/sample.png)

## What is it?

*Mosaic Video* is a project designed to transform live streaming videos into nice mosaics in realtime.

It uses ffmpeg to capture the video and position them in a grid.

## Features

* Generate mosaic from multiple video inputs;
* Audio support (first video input);
* Avaiable inputs: HLS.

## Roadmap

* Add support for RTMP inputs;
* Add support to save the mosaics into storage.

## Contributing

Feel free to contribute to this project by opening issues or pull requests.

`just test` to run the tests.

`just lint` to run linters.

`just deps` starts a mock API and Redis.

`just run` starts the go process to generate mosaics.

## Installation

We are still working on the installation process. For now, you can install it by cloning the repository and running `go run main.go`.

## Contributing

Feel free to contribute to this project by opening issues or pull requests.

`just test` to run the tests.

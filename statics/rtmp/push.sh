#!/bin/sh

ffmpeg -re -nostdin -i video.mp4 \
    -vcodec libx264 -preset:v ultrafast \
    -acodec aac \
    -f flv rtmp://localhost/mosaic-video/kansas
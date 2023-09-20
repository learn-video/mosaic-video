#!/bin/sh
ffmpeg \
    -i https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/level_2.m3u8 \
    -i https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/level_4.m3u8 \
    -filter_complex " \
        [0:v] setpts=PTS-STARTPTS, scale=qvga [a0]; \
        [1:v] setpts=PTS-STARTPTS, scale=qvga [a1]; \
        [a0][a1]xstack=inputs=2:layout=0_0|w0_0[out] \
        " \
    -map "[out]" \
    -c:v libx264 -t '30' -f matroska - | ffplay -autoexit -left 30 -top 30 -

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
    -c:v libx264 \
    -x264opts keyint=30:min-keyint=30:scenecut=-1 \
    -f hls \
    -hls_time 5 \
    -hls_start_number_source epoch \
    -hls_list_size 0 \
    -hls_segment_filename "output/segment%03d.ts" \
    output/playlist.m3u8

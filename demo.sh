#!/bin/sh

# ffmpeg \
#     -i https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/level_2.m3u8 \
#     -i https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/level_4.m3u8 \
#     -filter_complex " \
#         [0:v] setpts=PTS-STARTPTS, scale=qvga [a0]; \
#         [1:v] setpts=PTS-STARTPTS, scale=qvga [a1]; \
#         [a0][a1]xstack=inputs=2:layout=0_0|w0_0[out] \
#         " \
#     -map "[out]" \
#     -c:v libx264 \
#     -x264opts keyint=30:min-keyint=30:scenecut=-1 \
#     -f hls \
#     -hls_time 5 \
#     -hls_start_number_source epoch \
#     -hls_list_size 0 \
#     -hls_segment_filename "output/segment%03d.ts" \
#     output/playlist.m3u8

ffmpeg \
    -i https://rtvelivestream.akamaized.net/rtvesec/24h/24h_main_720.m3u8 \
    -i https://canadaremar2.todostreaming.es/live/peque-pequetv.m3u8 \
    -filter_complex " \
        [0:v] setpts=PTS-STARTPTS, scale=qvga [v0]; \
        [1:v] setpts=PTS-STARTPTS, scale=qvga [v1]; \
        [v0][v1]xstack=inputs=2:layout=0_0|w0_0[outv]; \
        [0:a] aresample=async=1 [a0]; \
        [1:a] aresample=async=1 [a1] \
        " \
    -map "[outv]" -c:v libx264 -b:v 800k -x264opts keyint=30:min-keyint=30:scenecut=-1 \
    -map "[a0]" -c:a aac -b:a 128k \
    -map "[a1]" -c:a aac -b:a 128k \
    -f hls \
    -hls_time 4 \
    -hls_list_size 6 \
    -hls_flags delete_segments \
    -hls_segment_filename "output/seg_%v_%03d.ts" \
    -var_stream_map "a:0,agroup:audio,default:yes,language:ENG a:1,agroup:audio,language:ENG v:0,agroup:audio" \
    -master_pl_name master.m3u8 \
    "output/playlist_%v.m3u8"


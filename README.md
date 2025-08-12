# Split

[Description](#description) . [Usage](#usage) . [Requirements](#requirements)

## Description

A little utility cooked up while learning go.

This will split a video file into many smaller files based on a .srt files
timestamps. This is useful for getting alot of voice training samples from a
single longer file. The resulting short ~8-12 second files can be used to train
TTS(text to speech) voices.

## Usage

./split -s <.srt file> -v <.wav file>

## Requirements

ffmpeg

#!/bin/bash

INPUTFILE=sample.wav
if [ $# -ge 1 ]; then
  INPUTFILE=$1
fi

d=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 $INPUTFILE)
file_duration=$(echo "$d/1" | bc)

echo "$INPUTFILE = ${file_duration} seconds"

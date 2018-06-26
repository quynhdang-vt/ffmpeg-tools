#!/bin/bash

if [ $# -ge 1 ]; then
  INPUT_FILE=$1
else
  echo "Please specify an audio file"
  exit 1
fi

ffprobe -show_entries stream=channels -of compact=p=0:nk=1 -v 0 $INPUT_FILE

## `ffprobe -v error -show_entries streams`
## shows a bunch of details as well

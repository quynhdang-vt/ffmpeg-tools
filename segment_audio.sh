#!/bin/bash

INPUTFILE=sample.wav
OUTPUT_FILE=output.txt
if [ $# -ge 1 ]; then
  INPUTFILE=$1
fi
if [ $# -ge 2 ]; then
  OUTPUT_FILE=$2
fi

echo silence_end: 0 > $OUTPUT_FILE
ffmpeg -hide_banner -i $INPUTFILE -af silencedetect=noise=-30dB:d=1 -f null - > /tmp/22.txt 2>&1
cat /tmp/22.txt | grep silencedetect | awk '{print $4,$5}'>> $OUTPUT_FILE
#cat $OUTPUT_FILE

## Output would be something like this
##
#silence_end: 0
#silence_start: -0.104
#silence_end: 6.4
#silence_start: 11.8
#silence_end: 16.256
#silence_start: 41.24
#silence_end: 43.008
#silence_start: 49.688
#silence_end: 56.064
#silence_start: 58.52

## to transform it so that it becoms silence_end, silence_start --> use this awk
   cat $OUTPUT_FILE | awk -F : '/^silence_end/{start_talking=$2}/^silence_start/{printf "%d,%d\n", start_talking*1000, $2*1000}'
## Then the output becomes this which is to indicate the non-silent segment
#
#0.000000,-0.104000
#6.400000,11.800000
#16.256000,41.240000
#43.008000,49.688000
#56.064000,58.520000

# from these output, we can start slicing the audio into
# audio_00000
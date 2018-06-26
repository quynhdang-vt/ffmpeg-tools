#!/bin/bash

OUTPUT_FILE=output.txt
if [ $# -ge 1 ]; then
  INPUT_FILE=$1
else
  echo "Please specify the audio file"
  exit 1
fi
if [ $# -ge 2 ]; then
  OUTPUT_FILE=$2
fi

d=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 $INPUT_FILE)
file_duration=$(echo "$d/1" | bc)

tmpfile1=/tmp/ttttt1.txt    #$(date "+%s")
tmpfile2=/tmp/yyyyy2.txt    #$(date "+%s")
#rm -f $tmpfile1

echo silence_end: 0 > $OUTPUT_FILE
ffmpeg -hide_banner -i $INPUT_FILE -af silencedetect=noise=-30dB:d=1 -f null - > $tmpfile1 2>&1
#ffmpeg -hide_banner -i $INPUT_FILE -af silencedetect=noise=-20dB:d=0.5 -f null - > $tmpfile1 2>&1
cat $tmpfile1 | grep silencedetect | awk '{print $4,$5}'>> $OUTPUT_FILE
echo silence_start: $file_duration >> $OUTPUT_FILE

rm -f $tmpfile2
cat $OUTPUT_FILE | awk -F : '/^silence_end/{start_talking=$2}/^silence_start/{printf "%d,%d\n", start_talking, $2}' >> $tmpfile2
mv $tmpfile2 $OUTPUT_FILE

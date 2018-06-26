#!/bin/bash

## format a time to HH:MM:SS format
function show_time() {
    local num=$1
    localmin=0
    local hour=0
    local day=0
    if((num>59));then
        ((sec=num%60))
        ((num=num/60))
        if((num>59));then
            ((min=num%60))
            ((num=num/60))
            if((num>23));then
                ((hour=num%24))
                ((day=num/24))
            else
                ((hour=num))
            fi
        else
            ((min=num))
        fi
    else
        ((sec=num))
    fi
    echo "$hour":"$min":"$sec"
}



#-------- sTART here -
#


INPUT_FILE=sample.wav
OUTPUT_FILE=output.txt
if [ $# -ge 1 ]; then
  INPUT_FILE=$1
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
#silence_start:xxx

#where xxx is of course the end of the file

## to transform it so that it becoms silence_end, silence_start --> use this awk
#rm -f $tmpfile1
rm -f $tmpfile2
cat $OUTPUT_FILE | awk -F : '/^silence_end/{start_talking=$2}/^silence_start/{printf "%d, %d\n", start_talking, $2}' >> $tmpfile2

mv $tmpfile2 $OUTPUT_FILE

#!/usr/bin/env bash

# Check if filename has been supplied as argument
if [ -z "$1" ]
  then
    echo "No filename supplied. Usage: ./script.sh filename"
    exit 1
fi

# Filename
FILENAME=$1

# Create a temporary backup of the original file
cp $FILENAME $FILENAME.bak

# Replace 8081 with 9081
sed -i 's/8081/9081/g' $FILENAME

echo "Replacement done in $FILENAME. Original file backed up as $FILENAME.bak."

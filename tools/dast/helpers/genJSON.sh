#!/bin/bash

# Generate a basic json file from parameters
# For each pair of parameters, set the first as key, and the second as value
# in the JSON file.

# Check if params are pairs
# For the json key vaule pairs
if [[ $(($#%2)) -ne 0 ]]; then
    echo "Illegal number of parameters" >&2
    exit 2
fi


echo "{"

# Generate json file
while test $# -gt 3
do
    echo "  \"$1\":\"$2\","
    shift
    shift
done

echo "  \"$1\":\"$2\""


echo "}"

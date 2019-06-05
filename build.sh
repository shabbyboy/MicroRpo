#!/bin/bash
set -e
read -p "Enter build path:" path

if [ -d $path ];
then
    read -p "Enter code path:" codepath
    if [ -d $codepath ];
    then
        dir=$(basename $codepath)
        cd $codepath
        #mkdir $path/$dir
        go build -o $path/$dir/$dir
        echo "build sucess"
    else
        echo "code path is not dir"
        exit 1
    fi
else
    echo "build path is not dir "
    exit 1
fi

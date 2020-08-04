#!/bin/bash

branch="$1"
image_tag='master'
organization='gravityhuborg'
image_name='gravity-data-extractor'

port=8098

case $branch in
    master) echo 'branch is master'; port=8071 ;;
    *) echo "Graceful exit..."; exit 0 ;;
esac

name="$image_name-$branch"

current_id=$(docker ps -a | grep "$name" | awk '{ print $1 }')

image_name="$organization/$image_name:$image_tag"

docker pull "$image_name"

docker stop "$current_id"
docker rm "$current_id"
docker run -itd -p "$port":80 --name "$name" "$image_name"

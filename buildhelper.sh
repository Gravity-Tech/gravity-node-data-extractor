#!/bin/sh

chmod +x ./buildcmd.sh
chmod +x ./buildpkg.sh

bash ./buildpkg.sh "$GOOS" "$GOARCH"
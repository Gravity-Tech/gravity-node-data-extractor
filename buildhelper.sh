#!/bin/sh

chmod +x ./buildcmd
chmod +x ./buildpkg
./buildcmd
./buildpkg "$GOOS" "$GOARCH"
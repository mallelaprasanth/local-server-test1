#!/bin/sh

# Mount GCS bucket
MOUNT_DIR=/go/tiles

mkdir $MOUNT_DIR
gcsfuse $INPUT_BUCKET $MOUNT_DIR

# Run server
./bin/mvt-server

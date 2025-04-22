#!/bin/bash

GO_POST_PROCESS_FILE="go fmt" openapi-generator generate \
  -i api/openapi.yaml \
    -g go-server \
    -o . \
    --enable-post-process-file \
    --additional-properties=packageName=api,sourceFolder=internal/api,apiPackage=api  \
    --git-repo-id=github.com/rotmanjanez/check24-gendev-7 \
    --git-user-id=rotmanjanez \

for f in $(find internal/api -name "model_*.go"); do
  modelname=$(echo $(basename $f) | sed 's/model_//; s/\.go//')
  echo "Renaming $f to pkg/models/$modelname.go"
  sed 's/package api/package models/' $f > pkg/models/$modelname.go
  rm $f
done

echo "Generated API code, manual adjustments are likely needed."
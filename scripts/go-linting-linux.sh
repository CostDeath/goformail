#!/bin/sh

# Go project root location relative to execution location
ROOT_LOCATION="."
# Linter executable location relative to Go project root.
EXEC_LOCATION="../../scripts/golangci-lint"

if command -v go
then
  cd $ROOT_LOCATION
  $EXEC_LOCATION run --out-format tab
else
  echo "Please install GoLang on your machine!"
fi